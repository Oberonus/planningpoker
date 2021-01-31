import axios from 'axios'

export default {
    async create() {
        const resp = await axios.post("games")
        return resp.data.game_id
    },

    async state(gameID) {
        const resp = await axios.get(`games/${gameID}`)
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