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
    },

    listenGame(gameID, callback) {
        this.socket.emit("join", gameID, res => {
            if (res !== 'ok') {
                this.listenStatus = this.STATUS_JOIN_FAILED
            } else {
                this.listenStatus = this.STATUS_JOINED
            }
        })
        this.socket.on("gameState", state => callback(state))
        this.listensGame = {gameID: gameID, callback: callback}
    },

    leaveGame() {
        this.socket.emit("leave")
        this.listensGame = null
        this.listenStatus = null
    }
}