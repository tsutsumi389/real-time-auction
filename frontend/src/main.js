import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { useBidderAuthStore } from './stores/bidderAuthStore'
import './assets/main.css'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)
app.use(router)

// アプリケーション起動時にログイン状態を復元（管理者と入札者の両方）
const authStore = useAuthStore()
const bidderAuthStore = useBidderAuthStore()

Promise.all([
  authStore.restoreUser(),
  bidderAuthStore.restoreUser()
]).then(() => {
  // ログイン状態復元後にアプリケーションをマウント
  app.mount('#app')
})
