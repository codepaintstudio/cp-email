import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import '@/styles/newmainpage.scss'
import '@/styles/mobile.scss'
import 'virtual:svg-icons-register'
import SvgIcon from './components/Svglcon.vue'
import router from './router'
import 'element-plus/theme-chalk/index.css'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App)
app.component('svg-icon', SvgIcon)
app.use(pinia)
app.use(router)
app.mount('#app')
