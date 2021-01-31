import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import user from '@/models/user'
import ls from 'local-storage'

user.name = ls.get('user_name')
user.id = ls.get('user_id')

axios.interceptors.request.use(req => {
    req.headers.authorization = `Bearer ${user.id}`;
    return req;
});
axios.defaults.baseURL = '/api/v1/'

Vue.config.productionTip = false

new Vue({
    router,
    vuetify,
    render: h => h(App)
}).$mount('#app')
