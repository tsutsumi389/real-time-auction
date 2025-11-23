<template>
  <div>
    <!-- モバイル用オーバーレイ -->
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-gray-600 bg-opacity-75 z-20 lg:hidden"
      @click="$emit('close')"
    ></div>

    <!-- サイドバー -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-30 w-64 bg-white border-r border-gray-200 transform transition-transform duration-200 ease-in-out lg:translate-x-0',
        isOpen ? 'translate-x-0' : '-translate-x-full',
      ]"
    >
      <div class="h-full flex flex-col">
        <!-- ロゴエリア（モバイル用） -->
        <div class="flex items-center justify-between h-16 px-4 border-b border-gray-200 lg:hidden">
          <h2 class="text-lg font-bold text-gray-900">メニュー</h2>
          <button
            @click="$emit('close')"
            class="p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
            aria-label="メニューを閉じる"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- ナビゲーションメニュー -->
        <nav class="flex-1 px-4 py-4 space-y-1 overflow-y-auto">
          <!-- ダッシュボード -->
          <router-link
            to="/admin/dashboard"
            :class="[
              'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
              isActive('/admin/dashboard')
                ? 'bg-blue-50 text-blue-700'
                : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
            ]"
            @click="handleNavigation"
          >
            <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              />
            </svg>
            ダッシュボード
          </router-link>

          <!-- 管理者管理（システム管理者のみ） -->
          <div v-if="authStore.isSystemAdmin">
            <div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">システム管理</div>
            <router-link
              to="/admin/admins"
              :class="[
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                isActive('/admin/admins')
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
              ]"
              @click="handleNavigation"
            >
              <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                />
              </svg>
              管理者一覧
            </router-link>
          </div>

          <!-- オークション管理 -->
          <div class="pt-4">
            <div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">オークション</div>
            <router-link
              to="/admin/auctions"
              :class="[
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                isActive('/admin/auctions')
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
              ]"
              @click="handleNavigation"
            >
              <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
                />
              </svg>
              オークション一覧
            </router-link>

            <router-link
              to="/admin/items"
              :class="[
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                isActive('/admin/items')
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
              ]"
              @click="handleNavigation"
            >
              <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                />
              </svg>
              商品一覧
              <span class="ml-auto text-xs bg-gray-200 text-gray-600 px-2 py-0.5 rounded">未実装</span>
            </router-link>
          </div>

          <!-- 入札者管理 -->
          <div class="pt-4">
            <div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">入札者管理</div>
            <router-link
              to="/admin/bidders"
              :class="[
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                isActive('/admin/bidders')
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
              ]"
              @click="handleNavigation"
            >
              <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
              入札者一覧
            </router-link>

            <router-link
              to="/admin/points"
              :class="[
                'flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors',
                isActive('/admin/points')
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900',
              ]"
              @click="handleNavigation"
            >
              <svg class="mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              ポイント管理
              <span class="ml-auto text-xs bg-gray-200 text-gray-600 px-2 py-0.5 rounded">未実装</span>
            </router-link>
          </div>
        </nav>

        <!-- フッター -->
        <div class="p-4 border-t border-gray-200">
          <div class="text-xs text-gray-500 text-center">
            <p>Real-time Auction System</p>
            <p class="mt-1">v1.0.0</p>
          </div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
})

const emit = defineEmits(['close'])

const route = useRoute()
const authStore = useAuthStore()

// アクティブ状態の判定
function isActive(path) {
  return route.path.startsWith(path)
}

// ナビゲーション時の処理（モバイルでサイドバーを閉じる）
function handleNavigation() {
  emit('close')
}
</script>
