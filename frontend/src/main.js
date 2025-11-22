import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import './assets/main.css'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)
app.use(router)

// アプリケーション起動時にログイン状態を復元
const authStore = useAuthStore()
authStore.restoreUser().then(() => {
  // ログイン状態復元後にアプリケーションをマウント
  app.mount('#app')
})
