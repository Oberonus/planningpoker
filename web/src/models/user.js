const axios = require('axios');
import ls from 'local-storage'

let _user = null;

export default class User {
    constructor() {
        if (_user) { return _user }
        this._name = ls.get('user_name') || "";
        this.id = ls.get('user_id') || "";
        this.registered = !!this._name;
        _user = this;
    }

    get identified() {
        return !!this._name && this.registered;
    }

    get name() { return this._name }
    set name(value) {
        if (!value || value === this._name) { return; }

        this.identified && this.update(value);
        this._name = value;
        ls.set('user_name', value);
    }

    authenticate() {
        if (!this.name) {
            throw "unknown user name, can not authenticate"
        }

        return axios.get("me")
            .then(rsp => {
                this._name = rsp.data.name;
                this.registered = true;
            })
            .catch(e => {
                if (e.response.status === 401) {
                    axios.post("register", {name: this._name})
                        .then(rsp => {
                            this.id = rsp.data.user_id
                            ls.set('user_id', this.id);
                            this.registered = true;
                        })
                } else {
                    console.error('request error: ', e);
                }
            })
    }

    update(name) {
        return axios
            .put("me", {name: name})
            .then(() => name);
    }
}
