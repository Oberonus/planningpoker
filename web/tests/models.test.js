const assert = require('node:assert/strict')
const fs = require('node:fs')
const path = require('node:path')
const vm = require('node:vm')
const { EventEmitter } = require('node:events')
const { test } = require('node:test')
const { transformSync } = require('@babel/core')
const Vue = require('vue')

// Use the frontend's existing Babel tooling to load alias-based modules with fake ports.
function loadModel(file, dependencies = {}, globals = {}) {
    const filename = path.join(__dirname, '../src', file)
    const content = fs.readFileSync(filename, 'utf8')
    const script = file.endsWith('.vue') ? content.match(/<script>([\s\S]*?)<\/script>/)[1] : content
    const source = transformSync(script, {
        filename,
        configFile: false,
        babelrc: false,
        plugins: [
            '@babel/plugin-transform-modules-commonjs',
            '@babel/plugin-proposal-class-properties',
        ],
    }).code
    const module = {exports: {}}

    vm.runInNewContext(source, {
        module,
        exports: module.exports,
        require: name => {
            assert.ok(name in dependencies, `Unexpected dependency: ${name}`)
            return dependencies[name]
        },
        setTimeout,
        clearTimeout,
        URL,
        ...globals,
    }, {filename})

    return module.exports
}

function socketFixture() {
    const events = new EventEmitter()
    const requests = []
    const socket = {
        connected: true,
        on: events.on.bind(events),
        once: events.once.bind(events),
        off: events.off.bind(events),
        emit: (...args) => requests.push(args),
    }

    return {socket, events, requests}
}

function timerFixture() {
    const timers = new Set()

    return {
        timers,
        setTimeout: callback => {
            timers.add(callback)
            return callback
        },
        clearTimeout: callback => timers.delete(callback),
        expire: () => {
            for (const callback of [...timers]) {
                timers.delete(callback)
                callback()
            }
        },
    }
}

test('actions confirm success and reject server errors instead of silently succeeding', async () => {
    const {socket, requests} = socketFixture()
    const timers = timerFixture()
    const game = loadModel('models/game.js', {'@/notifier/notifier': {socket}}, timers).default

    const creation = game.create('Planning', '', game.decks[0])
    requests[0][2]({game_id: 'new-game'})
    assert.equal(await creation, 'new-game')
    assert.equal(timers.timers.size, 0)

    const reveal = game.reveal()
    const rejection = assert.rejects(reveal, /Game request failed/)
    requests[1][1]('Unable to reveal')
    await rejection
    assert.equal(timers.timers.size, 0)
})

test('timed-out disconnected creation is not sent when the connection eventually returns', async () => {
    const {socket, events, requests} = socketFixture()
    const timers = timerFixture()
    socket.connected = false
    const game = loadModel('models/game.js', {'@/notifier/notifier': {socket}}, timers).default

    const creation = game.create('Planning', '', game.decks[0])
    const rejection = assert.rejects(creation, /timed out/)
    assert.equal(requests.length, 0)
    timers.expire()
    await rejection

    socket.connected = true
    events.emit('connect')
    assert.equal(requests.length, 0)
    assert.equal(events.listenerCount('connect'), 0)
})

test('unacknowledged actions time out and unavailable connections report failure', async () => {
    const {socket} = socketFixture()
    const timers = timerFixture()
    const game = loadModel('models/game.js', {'@/notifier/notifier': {socket}}, timers).default

    const restart = game.restart()
    const rejection = assert.rejects(restart, /timed out/)
    timers.expire()
    await rejection

    const disconnected = loadModel('models/game.js', {'@/notifier/notifier': {socket: null}}, timers).default
    await assert.rejects(disconnected.restart(), /unavailable/)
})

test('selection and confidence follow server state, with Unicode votes encoded once', async () => {
    const {socket, requests} = socketFixture()
    const timers = timerFixture()
    const game = loadModel('models/game.js', {'@/notifier/notifier': {socket}}, timers).default
    let receiveState
    const State = loadModel('models/state.js', {
        '@/models/game': game,
        '@/notifier/notifier': {
            listenGame: (id, callback) => { receiveState = callback },
        },
    }).default
    const state = new State('game')

    const vote = state.vote('½')
    assert.equal(state.voted(), false)
    assert.equal(requests[0][1].vote, '%C2%BD')
    requests[0][2]('')
    await vote
    assert.equal(state.voted(), false)

    receiveState({voted_card: '½', confidence: 'normal', players: []})
    assert.equal(state.isActive('½'), true)

    const confidence = state.changeConfidence('high')
    assert.equal(requests[1][1].vote, '%C2%BD')
    assert.equal(requests[1][1].confidence, 'high')
    assert.equal(state.confidence, 'normal')
    const rejection = assert.rejects(confidence, /failed/)
    requests[1][2]('Unable to vote')
    await rejection
    assert.equal(state.confidence, 'normal')

    const unvote = state.vote('½')
    assert.equal(requests[2][0], 'unvote')
    requests[2][1]('ok')
    await unvote
    receiveState({voted_card: '', confidence: '', players: []})
    assert.equal(state.voted(), false)
})

test('rejoining replaces state listeners and ignores acknowledgements from earlier games', () => {
    const {socket, events, requests} = socketFixture()
    const timers = timerFixture()
    const notifier = loadModel('notifier/notifier.js', {'socket.io-client': () => socket}, timers).default
    notifier.socket = socket
    let firstUpdates = 0
    let currentUpdates = 0

    notifier.listenGame('first', () => { firstUpdates++ })
    notifier.listenGame('current', () => { currentUpdates++ })
    assert.equal(events.listenerCount('gameState'), 1)
    assert.equal(timers.timers.size, 1)

    requests[0][2]('failure')
    assert.equal(notifier.listenStatus, null)
    requests[1][2]('ok')
    assert.equal(notifier.listenStatus, notifier.STATUS_JOINED)
    events.emit('gameState', {})
    assert.equal(firstUpdates, 0)
    assert.equal(currentUpdates, 1)

    notifier.listenGame('current', () => { currentUpdates++ })
    timers.expire()
    assert.equal(notifier.listenStatus, notifier.STATUS_JOIN_FAILED)

    notifier.leaveGame()
    assert.equal(events.listenerCount('gameState'), 0)
    requests[2][2]('ok')
    assert.equal(notifier.listenStatus, null)
})

test('ticket validation accepts web links and rejects executable and incomplete links', () => {
    const {ticketURLRules} = loadModel('models/validation.js')
    const validate = ticketURLRules[0]

    assert.equal(validate(''), true)
    assert.equal(validate('  https://example.com/issue/42  '), true)
    assert.equal(validate('http://localhost:18080/ticket'), true)
    assert.notEqual(validate('javascript:alert(1)'), true)
    assert.notEqual(validate('example.com/issue'), true)
})

test('clipboard fallback reports success only when copying succeeds, and dismissal settles the dialog', async () => {
    let canCopy = false
    const input = {focus: () => {}, select: () => {}}
    const InviteDialog = loadModel('components/InviteDialog.vue', {}, {
        document: {execCommand: () => canCopy},
        requestAnimationFrame: callback => callback(),
    }).default
    const dialog = new Vue(InviteDialog)
    dialog.$refs.textToCopy = {$el: {querySelector: () => input}}

    const copying = dialog.open('http://localhost/game')
    await Vue.nextTick()
    dialog.copy()
    assert.equal(dialog.copyFailed, true)
    assert.equal(dialog.show, true)

    canCopy = true
    dialog.copy()
    assert.equal(await copying, true)
    await Vue.nextTick()

    const dismissed = dialog.open('http://localhost/game')
    await Vue.nextTick()
    assert.equal(dialog.copyFailed, false)
    dialog.show = false
    assert.equal(await dismissed, false)
    dialog.$destroy()
})
