import io from 'socket.io-client';

export default {
    STATUS_CONNECTED: "connected",
    STATUS_RECONNECTING: "reconnecting",
    STATUS_JOIN_FAILED: "join-failed",
    STATUS_JOINED: "joined",

    socket: null,
    listensGame: null,
    status: null,
    listenStatus: null,
    stateHandler: null,
    joinTimeout: null,
    joinSequence: 0,

    connect(token) {
        // we should create only one socket per session
        if (this.socket != null) {
            return
        }

        this.socket = io("/", {
            query: {
                token: token,
            },
            transports: ["websocket"],
        })

        this.socket.on('connect', () => {
            this.status = this.STATUS_CONNECTED
        });

        this.socket.on('reconnect', () => {
            this.status = this.STATUS_CONNECTED
            if (this.listensGame != null) {
                this.listenGame(this.listensGame.gameID, this.listensGame.callback)
            }
        })

        this.socket.on('reconnecting', () => {
            this.status = this.STATUS_RECONNECTING
        })

        this.socket.on('disconnect', () => {
            this.status = this.STATUS_RECONNECTING
        })

        this.socket.on('connect_error', () => {
            this.status = this.STATUS_RECONNECTING
        })
    },

    listenGame(gameID, callback) {
        clearTimeout(this.joinTimeout)

        if (this.stateHandler) {
            this.socket.off('gameState', this.stateHandler)
        }

        this.listenStatus = null
        this.stateHandler = state => callback(state)
        this.socket.on('gameState', this.stateHandler)
        this.listensGame = {gameID: gameID, callback: callback}

        const sequence = ++this.joinSequence
        this.joinTimeout = setTimeout(() => {
            this.listenStatus = this.STATUS_JOIN_FAILED
        }, 10000)

        this.socket.emit("join", gameID, res => {
            if (sequence !== this.joinSequence) {
                return
            }

            clearTimeout(this.joinTimeout)

            if (res !== 'ok') {
                this.listenStatus = this.STATUS_JOIN_FAILED
            } else {
                this.listenStatus = this.STATUS_JOINED
            }
        })
    },

    leaveGame() {
        clearTimeout(this.joinTimeout)
        this.joinSequence++

        if (this.stateHandler) {
            this.socket.off('gameState', this.stateHandler)
            this.stateHandler = null
        }

        this.socket.emit("leave")
        this.listensGame = null
        this.listenStatus = null
    }
}
