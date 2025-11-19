<template>
  <div id="app">
    <header class="header">
      <h1>🏇 Real-Time Auction System</h1>
      <p>競走馬セリをモデルとしたリアルタイムオークションシステム</p>
    </header>
    
    <main class="main">
      <router-view />
    </main>

    <footer class="footer">
      <p>© 2025 Real-Time Auction System</p>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  console.log('Real-Time Auction System - Frontend')
  console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL)
  console.log('WebSocket URL:', import.meta.env.VITE_WS_URL)
})
</script>

<style scoped>
.header {
  padding: 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-align: center;
}

.header h1 {
  margin: 0;
  font-size: 2rem;
}

.header p {
  margin: 0.5rem 0 0;
  opacity: 0.9;
}

.main {
  min-height: calc(100vh - 200px);
  padding: 2rem;
}

.footer {
  padding: 1rem;
  background: #f5f5f5;
  text-align: center;
  color: #666;
}
</style>
