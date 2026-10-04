const assert = require('node:assert/strict')
const { test } = require('node:test')
const io = require('socket.io-client')

const baseURL = process.env.POKER_TEST_URL
assert.ok(baseURL, 'POKER_TEST_URL must point to a disposable test service')
const timeout = 5000

async function request(method, path, body, token) {
    const headers = {'Content-Type': 'application/json'}
    if (token) headers.Authorization = `Bearer ${token}`

    const response = await fetch(`${baseURL}${path}`, {
        method,
        headers,
        body: method === 'GET' || body === undefined ? undefined : JSON.stringify(body),
        signal: AbortSignal.timeout(timeout),
    })

    return {status: response.status, body: await response.json()}
}

async function register(name) {
    const response = await request('POST', '/api/v1/register', {name})
    assert.equal(response.status, 200)
    assert.ok(response.body.user_id)
    return response.body.user_id
}

function connect(t, token) {
    const socket = io(baseURL, {
        query: {token},
        transports: ['websocket'],
        reconnection: false,
        forceNew: true,
        autoConnect: false,
        timeout,
    })
    t.after(() => socket.disconnect())

    return new Promise((resolve, reject) => {
        const timer = setTimeout(() => reject(new Error('Socket connection timed out')), timeout)
        socket.once('connect', () => {
            clearTimeout(timer)
            resolve(socket)
        })
        socket.once('connect_error', error => {
            clearTimeout(timer)
            reject(error)
        })
        socket.connect()
    })
}

function emit(socket, event, payload) {
    return new Promise((resolve, reject) => {
        const timer = setTimeout(() => reject(new Error(`${event} acknowledgement timed out`)), timeout)
        const acknowledge = response => {
            clearTimeout(timer)
            resolve(response)
        }

        if (payload === undefined) socket.emit(event, acknowledge)
        else socket.emit(event, payload, acknowledge)
    })
}

// Subscribe before the action: state publication can precede its acknowledgement.
function nextState(socket, predicate) {
    return new Promise((resolve, reject) => {
        const timer = setTimeout(() => {
            socket.off('gameState', receive)
            reject(new Error('Expected game state was not published'))
        }, timeout)

        function receive(state) {
            if (!predicate(state)) return
            clearTimeout(timer)
            socket.off('gameState', receive)
            resolve(state)
        }

        socket.on('gameState', receive)
    })
}

async function create(socket, everyoneCanReveal = false) {
    const response = await emit(socket, 'create', {
        name: 'Checkout estimate',
        url: 'https://example.com/TICKET-1',
        cards_deck: {name: 'Custom', types: ['1', '5', '½', '?']},
        everyone_can_reveal: everyoneCanReveal,
    })
    assert.ok(response.game_id)
    return response.game_id
}

async function join(socket, gameID, count) {
    const state = nextState(socket, state => state.players.length === count)
    assert.equal(await emit(socket, 'join', gameID), 'ok')
    return state
}

function player(state, name) {
    const result = state.players.find(player => player.name === name)
    assert.ok(result, `Missing player ${name}`)
    return result
}

test('identity: register, authenticate, rename and reject invalid requests', {timeout: 20000}, async () => {
    assert.equal((await request('GET', '/alive')).status, 200)
    const token = await register('Alice')
    assert.deepEqual(await request('GET', '/api/v1/me', undefined, token), {
        status: 200,
        body: {name: 'Alice'},
    })

    assert.equal((await request('PUT', '/api/v1/me', {name: 'Alicia'}, token)).status, 200)
    assert.equal((await request('GET', '/api/v1/me', undefined, token)).body.name, 'Alicia')

    for (const method of ['GET', 'PUT']) {
        for (const invalidToken of [undefined, 'unknown-user']) {
            assert.deepEqual(await request(method, '/api/v1/me', {name: 'Intruder'}, invalidToken), {
                status: 401,
                body: {error: 'unauthorized'},
            })
        }
    }

    const response = await fetch(`${baseURL}/api/v1/register`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: '{invalid',
        signal: AbortSignal.timeout(timeout),
    })
    assert.equal(response.status, 400)
    assert.equal((await request('GET', '/api/v1/me', undefined, token)).body.name, 'Alicia')
})

test('game round: private votes, live updates, reveal permissions and restart', {timeout: 30000}, async t => {
    const aliceToken = await register('Alice')
    const alice = await connect(t, aliceToken)
    const bob = await connect(t, await register('Bob'))
    const gameID = await create(alice)
    const initial = await join(alice, gameID, 1)
    assert.equal(initial.name, 'Checkout estimate')
    assert.equal(initial.ticket_url, 'https://example.com/TICKET-1')
    assert.deepEqual(initial.cards_deck, {name: 'Custom', cards: ['1', '5', '½', '?']})
    assert.equal(initial.can_reveal, true)

    const bobState = await join(bob, gameID, 2)
    assert.equal(bobState.can_reveal, false)
    assert.equal(await emit(bob, 'reveal'), 'error')

    const aliceVote = nextState(alice, state => state.voted_card === '½')
    const bobView = nextState(bob, state => player(state, 'Alice').voted_card === '*')
    assert.equal(await emit(alice, 'vote', {vote: encodeURIComponent('½'), confidence: 'high'}), '')
    assert.equal((await aliceVote).confidence, 'high')
    const hidden = await bobView
    assert.equal(hidden.state, 'started')
    assert.equal(player(hidden, 'Alice').confidence, '')

    assert.equal(await emit(bob, 'vote', {vote: '99', confidence: 'normal'}), 'error')
    const bobVote = nextState(bob, state => state.voted_card === '5')
    assert.equal(await emit(bob, 'vote', {vote: '5', confidence: 'normal'}), '')
    await bobVote

    const removed = nextState(bob, state => state.voted_card === '')
    assert.equal(await emit(bob, 'unvote'), 'ok')
    assert.equal(player(await removed, 'Bob').voted_card, '')

    const renamed = nextState(bob, state => state.players.some(player => player.name === 'Alicia'))
    assert.equal((await request('PUT', '/api/v1/me', {name: 'Alicia'}, aliceToken)).status, 200)
    assert.equal(player(await renamed, 'Alicia').voted_card, '*')

    const updated = nextState(bob, state => state.name === 'Next estimate')
    assert.equal(await emit(alice, 'update', {name: 'Next estimate', ticket_url: 'https://example.com/TICKET-2'}), 'ok')
    assert.equal((await updated).ticket_url, 'https://example.com/TICKET-2')

    const revealed = nextState(bob, state => state.state === 'finished')
    assert.equal(await emit(alice, 'reveal'), 'ok')
    const finished = await revealed
    assert.equal(player(finished, 'Alicia').voted_card, '½')
    assert.equal(player(finished, 'Alicia').confidence, 'high')
    assert.equal(await emit(bob, 'vote', {vote: '5', confidence: 'normal'}), 'error')
    assert.equal(await emit(alice, 'unvote'), 'error')

    const restarted = nextState(bob, state => state.state === 'started' && state.players.every(player => player.voted_card === ''))
    assert.equal(await emit(alice, 'restart'), 'ok')
    assert.equal((await restarted).voted_card, '')

    const left = nextState(alice, state => state.players.length === 1)
    assert.equal(await emit(bob, 'leave'), 'ok')
    assert.equal(player(await left, 'Alicia').voted_card, '')
})

test('reconnect restores votes; a departed voter remains until the next round', {timeout: 20000}, async t => {
    const token = await register('Owner')
    const owner = await connect(t, token)
    const guest = await connect(t, await register('Guest'))
    const gameID = await create(owner, true)
    await join(owner, gameID, 1)
    assert.equal((await join(guest, gameID, 2)).can_reveal, true)

    const voted = nextState(owner, state => state.voted_card === '5')
    assert.equal(await emit(owner, 'vote', {vote: '5', confidence: 'low'}), '')
    await voted
    owner.disconnect()

    const reconnected = await connect(t, token)
    const restored = await join(reconnected, gameID, 2)
    assert.equal(restored.voted_card, '5')
    assert.equal(restored.confidence, 'low')

    assert.equal(await emit(reconnected, 'leave'), 'ok')
    const revealed = nextState(guest, state => state.state === 'finished')
    assert.equal(await emit(guest, 'reveal'), 'ok')
    assert.equal(player(await revealed, 'Owner').voted_card, '5')

    const restarted = nextState(guest, state => state.state === 'started' && state.players.length === 1)
    assert.equal(await emit(guest, 'restart'), 'ok')
    assert.equal(player(await restarted, 'Guest').voted_card, '')
})

test('invalid game requests are rejected and the connection stays usable', {timeout: 15000}, async t => {
    const socket = await connect(t, await register('Player'))
    assert.equal(await emit(socket, 'join', 'missing-game'), 'error')
    assert.equal(await emit(socket, 'vote', {vote: '%broken', confidence: 'normal'}), 'error')
    assert.equal(await emit(socket, 'create', {
        name: 'Invalid deck',
        cards_deck: {name: 'Empty', types: []},
    }), 'error')

    const gameID = await create(socket)
    assert.equal((await join(socket, gameID, 1)).state, 'started')
})

test('unregistered clients cannot keep a game connection', {timeout: 10000}, async t => {
    const socket = io(baseURL, {
        query: {token: 'unknown-user'},
        transports: ['websocket'],
        reconnection: false,
        forceNew: true,
        autoConnect: false,
        timeout,
    })
    t.after(() => socket.disconnect())
    socket.on('gameState', () => assert.fail('Unauthenticated client received game state'))

    await new Promise((resolve, reject) => {
        const timer = setTimeout(() => reject(new Error('Unauthenticated connection stayed open')), timeout)
        const rejected = () => {
            clearTimeout(timer)
            resolve()
        }
        socket.once('connect_error', rejected)
        socket.once('disconnect', rejected)
        socket.connect()
    })
    assert.equal(socket.connected, false)
})

test('simultaneous player votes both survive and reach subscribers', {timeout: 15000}, async t => {
    const alice = await connect(t, await register('Concurrent Alice'))
    const bob = await connect(t, await register('Concurrent Bob'))
    const gameID = await create(alice)
    await join(alice, gameID, 1)
    await join(bob, gameID, 2)

    const bothVoted = state => state.players.length === 2 && state.players.every(player => player.voted_card === '*')
    const aliceState = nextState(alice, bothVoted)
    const bobState = nextState(bob, bothVoted)
    assert.deepEqual(await Promise.all([
        emit(alice, 'vote', {vote: '1', confidence: 'high'}),
        emit(bob, 'vote', {vote: '5', confidence: 'low'}),
    ]), ['', ''])
    assert.equal((await aliceState).voted_card, '1')
    assert.equal((await bobState).voted_card, '5')

    const revealed = nextState(bob, state => state.state === 'finished')
    assert.equal(await emit(alice, 'reveal'), 'ok')
    const state = await revealed
    assert.equal(player(state, 'Concurrent Alice').voted_card, '1')
    assert.equal(player(state, 'Concurrent Bob').voted_card, '5')
})
