import notifier from "@/notifier/notifier";

export default {
    decks: [
        {name: "T-Shirt", types: ["XXS", "XS", "S", "M", "L", "XL", "XXL", "?"]},
        {name: "Fibonacci", types: ["0", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "?"]},
        {name: "Custom fibonacci", types: ["0", "½", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?"]},
        {name: "Powers of 2", types: ["0", "1", "2", "4", "8", "16", "32", "64", "?"]},
    ],

    create(name, url, deck) {
        const ob = {
            name: name,
            url: url,
            cards_deck: deck,
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

    async vote(gameID, vote, confidence) {
        notifier.socket.emit("vote", {
            vote: vote,
            confidence: confidence,
        })
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