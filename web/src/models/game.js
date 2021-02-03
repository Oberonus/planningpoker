import axios from 'axios'

export default {
    async create() {
        const resp = await axios.post("games")
        return resp.data.game_id
    },

    state(gameID) {
        return axios
            .get(`games/${gameID}`)
            .then(resp => resp.data);
    },

    join(gameID) {
        return axios.post(`games/${gameID}/join`)
            .then((resp) => {
                return resp.data;
            })
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
