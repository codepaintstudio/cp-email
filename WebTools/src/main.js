import { createApp } from 'vue'
import App from './App.vue'
import '@vueup/vue-quill/dist/vue-quill.snow.css'
import { createPinia } from 'pinia'
import piniaPersist from 'pinia-plugin-persist'
import '@/styles/newmainpage.scss'
import '@/styles/mobile.scss'
import 'virtual:svg-icons-register'
import SvgIcon from './components/Svglcon.vue'
import router from './router'
import 'element-plus/theme-chalk/index.css'

const pinia = createPinia()
pinia.use(piniaPersist)

const app = createApp(App)
app.component('svg-icon', SvgIcon)
app.use(pinia)
app.use(router)
app.mount('#app')
