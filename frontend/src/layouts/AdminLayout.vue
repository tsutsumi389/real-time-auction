<template>
  <div class="min-h-screen bg-gray-50">
    <!-- サイドバー -->
    <AdminSidebar :is-open="sidebarOpen" @close="closeSidebar" />

    <!-- メインコンテンツエリア -->
    <div class="lg:pl-64">
      <!-- ヘッダー -->
      <AdminHeader @toggle-sidebar="toggleSidebar" />

      <!-- ページコンテンツ -->
      <main class="flex-1">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'

const sidebarOpen = ref(false)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

function closeSidebar() {
  sidebarOpen.value = false
}
</script>
