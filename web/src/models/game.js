import notifier from "@/notifier/notifier"

// A bounded acknowledgement keeps failed actions from leaving the UI busy indefinitely.
function request(event, payload) {
    return new Promise((resolve, reject) => {
        const socket = notifier.socket
        if (!socket) {
            reject(new Error('Game connection is unavailable'))
            return
        }

        const timeout = setTimeout(() => {
            socket.off('connect', sendRequest)
            reject(new Error('Game request timed out'))
        }, 10000)
        const acknowledge = response => {
            clearTimeout(timeout)

            if (response === 'ok' || response === '' || (response && response.game_id)) {
                resolve(response)
                return
            }

            reject(new Error('Game request failed'))
        }

        function sendRequest() {
            if (payload === undefined) {
                socket.emit(event, acknowledge)
            } else {
                socket.emit(event, payload, acknowledge)
            }
        }

        if (socket.connected) {
            sendRequest()
        } else {
            // Do not put timed-out actions in Socket.IO's disconnected send buffer.
            socket.once('connect', sendRequest)
        }
    })
}

export default {
    decks: [
        {name: "T-Shirt", types: ["XXS", "XS", "S", "M", "L", "XL", "XXL", "?"]},
        {name: "Fibonacci", types: ["0", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "?"]},
        {name: "Custom fibonacci", types: ["0", "½", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?"]},
        {name: "Powers of 2", types: ["0", "1", "2", "4", "8", "16", "32", "64", "?"]},
    ],

    async create(name, url, deck) {
        const response = await request('create', {
            name,
            url,
            cards_deck: deck,
            everyone_can_reveal: true,
        })

        return response.game_id
    },

    update(gameID, name, url) {
        return request('update', {
            name,
            ticket_url: url,
        })
    },

    leave() {
        return request('leave')
    },

    vote(gameID, vote, confidence) {
        return request('vote', {
            vote: encodeURIComponent(vote),
            confidence,
        })
    },

    reveal() {
        return request('reveal')
    },

    unVote() {
        return request('unvote')
    },

    restart() {
        return request('restart')
    },
}
