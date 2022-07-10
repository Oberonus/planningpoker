import notifier from "@/notifier/notifier";

const cardsDeck = {
    name: "T-Shirt",
    types: ["XXS", "XS", "S", "M", "L", "XL", "XXL", "?"]
}

export default {
    create(name, url) {
        const ob = {
            name: name,
            url: url,
            cards_deck: cardsDeck,
            everyone_can_reveal: true,
        }

        return new Promise((resolve) => {
            notifier.socket.emit("create", ob, (data) => {
                resolve(data.game_id)
            })
        });
    },

    async update(gameID, name, url) {
        notifier.socket.emit("update", {
            name: name,
            ticket_url: url
        })
    },

    async leave() {
        notifier.socket.emit("leave")
    },

    async vote(gameID, vote) {
        notifier.socket.emit("vote", vote)
    },

    async reveal() {
        notifier.socket.emit("reveal")
    },

    async unVote() {
        notifier.socket.emit("unvote")
    },

    async restart() {
        notifier.socket.emit("restart")
    },
}