import axios from 'axios'

const cardsDeck = {
    name: "T-Shirt",
    types: ["XXS", "XS", "S", "M", "L", "XL", "XXL", "?"]
}

export default {
    async create(name, url) {
        const resp = await axios.post("games", {
            name: name,
            url: url,
            cards_deck: cardsDeck,
            everyone_can_reveal: true,
        })
        return resp.data.game_id
    },

    async update(gameID, name, url) {
        const resp = await axios.put(`games/${gameID}`, {
            name: name,
            url: url,
        })
        return resp.data
    },

    async state(gameID, lastChangeID) {
        const resp = await axios.get(`games/${gameID}?lastChangeID=${lastChangeID}`)
        return resp.data
    },

    async join(gameID) {
        const resp = await axios.post(`games/${gameID}/join`)
        return resp.data
    },

    async vote(gameID, vote) {
        const resp = await axios.post(`games/${gameID}/votes/${vote}`)
        return resp.data
    },

    async reveal(gameID) {
        const resp = await axios.post(`games/${gameID}/reveal`)
        return resp.data
    },

    async unVote(gameID) {
        const resp = await axios.post(`games/${gameID}/unvote`)
        return resp.data
    },

    async restart(gameID) {
        const resp = await axios.post(`games/${gameID}/restart`)
        return resp.data
    },
}