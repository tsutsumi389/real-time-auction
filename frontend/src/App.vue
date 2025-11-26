<template>
  <div id="app">
    <!-- 管理画面レイアウト -->
    <component :is="layout">
      <router-view />
    </component>
    <!-- トースト通知コンテナ -->
    <ToastContainer />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const route = useRoute()

// 現在のルートのメタ情報に基づいてレイアウトを決定
const layout = computed(() => {
  // hideLayoutがtrueの場合はレイアウトなし
  if (route.meta.hideLayout) {
    return 'div'
  }
  // 管理画面の場合はAdminLayout
  if (route.path.startsWith('/admin') && route.name !== 'admin-login') {
    return AdminLayout
  }
  // その他はDefaultLayout
  return DefaultLayout
})

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
