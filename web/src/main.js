import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

axios.defaults.baseURL = "http://localhost:8080/"
axios.defaults.withCredentials = false
axios.defaults.headers = {
    // 'Content-Type': 'application/json; charset=UTF-8',
    // 'Access-Control-Allow-Origin': '*',
    // 'Access-Control-Allow-Headers': '*',
    // 'Accept': 'Token',
}
const app = createApp(App)
app.config.globalProperties.axios = axios
app.use(router).mount('#app')
