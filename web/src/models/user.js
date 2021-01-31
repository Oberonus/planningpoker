const axios = require('axios');
import ls from 'local-storage'

export default {
    name: "",
    id: "",

    identified() {
        console.log(`user name === '${this.name}'`)
        return this.name !== undefined && this.name !== "" && this.name != null
    },

    async authenticate() {
        if (!this.identified()) {
            throw "unknown user name, can not authenticate"
        }

        try {
            const rsp = await axios.get("me")
            // actualize name
            this.name = rsp.data.name
        } catch (e) {
            if (e.response.status === 401) {
                const rsp = await axios.post("register", {name: this.name})
                this.id = rsp.data.user_id
            } else {
                throw e
            }
        }

        ls.set('user_name', this.name)
        ls.set('user_id', this.id)
    },

    async update(name) {
        await axios.put("me", {name: name})
        this.name = name
    }
}