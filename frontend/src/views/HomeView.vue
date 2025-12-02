<template>
  <div class="container mx-auto py-8 px-4">
    <!-- Welcome Card -->
    <Card class="mb-8 p-8 text-center">
      <h2 class="text-3xl font-bold mb-2 text-primary">🎉 開発環境セットアップ完了</h2>
      <p class="text-muted-foreground">Docker + Vite + Vue.js 3 + Shadcn Vue で動作中</p>
    </Card>

    <!-- Status Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <!-- API Test Card -->
      <Card class="p-6">
        <h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
          📡 API接続テスト
        </h3>
        <Button
          class="w-full mb-4"
          :disabled="apiLoading"
          @click="testApi"
        >
          {{ apiLoading ? 'テスト中...' : 'APIをテスト' }}
        </Button>
        <div
          v-if="apiResult"
          :class="[
            'p-4 rounded-md text-sm',
            apiSuccess
              ? 'bg-green-50 text-green-800 border border-green-200'
              : 'bg-red-50 text-red-800 border border-red-200'
          ]"
        >
          {{ apiResult }}
        </div>
      </Card>

      <!-- WebSocket Test Card -->
      <Card class="p-6">
        <h3 class="text-xl font-semibold mb-4 flex items-center gap-2">
          🔌 WebSocket接続テスト
        </h3>
        <Button
          class="w-full mb-4"
          variant="secondary"
          :disabled="wsLoading"
          @click="testWebSocket"
        >
          {{ wsLoading ? 'テスト中...' : 'WebSocketをテスト' }}
        </Button>
        <div
          v-if="wsResult"
          :class="[
            'p-4 rounded-md text-sm',
            wsSuccess
              ? 'bg-green-50 text-green-800 border border-green-200'
              : 'bg-red-50 text-red-800 border border-red-200'
          ]"
        >
          {{ wsResult }}
        </div>
      </Card>
    </div>

    <!-- Design System Showcase -->
    <Card class="mb-8 p-6">
      <h3 class="text-xl font-semibold mb-4">🎨 Shadcn Vue + Tailwind CSS デモ</h3>
      <div class="space-y-4">
        <div class="flex flex-wrap gap-2">
          <Button>Default</Button>
          <Button variant="secondary">Secondary</Button>
          <Button variant="destructive">Destructive</Button>
          <Button variant="outline">Outline</Button>
          <Button variant="ghost">Ghost</Button>
          <Button variant="link">Link</Button>
        </div>
        <div class="flex flex-wrap gap-2">
          <Button size="sm">Small</Button>
          <Button>Default</Button>
          <Button size="lg">Large</Button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
          <div class="p-4 rounded-lg bg-primary text-primary-foreground">
            Primary Color
          </div>
          <div class="p-4 rounded-lg bg-secondary text-secondary-foreground">
            Secondary Color
          </div>
          <div class="p-4 rounded-lg bg-accent text-accent-foreground">
            Accent Color
          </div>
        </div>
      </div>
    </Card>

    <!-- Next Steps -->
    <Card class="p-6">
      <h3 class="text-xl font-semibold mb-4">📋 次のステップ</h3>
      <ul class="space-y-2 text-muted-foreground">
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>データベーススキーマの設計と実装</span>
        </li>
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>JWT認証システムの実装</span>
        </li>
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>ユーザー管理APIの実装</span>
        </li>
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>オークション管理機能の実装</span>
        </li>
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>WebSocketリアルタイム通信の実装</span>
        </li>
        <li class="flex items-start gap-2">
          <span class="text-primary">•</span>
          <span>フロントエンドUIコンポーネントの構築</span>
        </li>
      </ul>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const apiLoading = ref(false)
const apiResult = ref('')
const apiSuccess = ref(false)

const wsLoading = ref(false)
const wsResult = ref('')
const wsSuccess = ref(false)

// 動的にAPIとWebSocket URLを取得（ローカルネットワーク対応）
import { getApiBaseUrl, getWsUrl } from '@/config/api'
const API_BASE_URL = getApiBaseUrl()
const WS_URL = getWsUrl()

async function testApi() {
  apiLoading.value = true
  apiResult.value = ''

  try {
    const response = await axios.get(`${API_BASE_URL}/ping`)
    apiSuccess.value = true
    apiResult.value = `✓ 成功: ${JSON.stringify(response.data)}`
  } catch (error) {
    apiSuccess.value = false
    apiResult.value = `✗ エラー: ${error.message}`
  } finally {
    apiLoading.value = false
  }
}

function testWebSocket() {
  wsLoading.value = true
  wsResult.value = ''

  try {
    const ws = new WebSocket(WS_URL)

    ws.onopen = () => {
      wsSuccess.value = true
      wsResult.value = '✓ WebSocket接続成功'
      ws.close()
      wsLoading.value = false
    }

    ws.onerror = () => {
      wsSuccess.value = false
      wsResult.value = '✗ WebSocket接続エラー（まだ実装されていません）'
      wsLoading.value = false
    }

    ws.onclose = () => {
      if (!wsSuccess.value && !wsResult.value) {
        wsSuccess.value = false
        wsResult.value = '✗ WebSocket接続エラー'
      }
      wsLoading.value = false
    }
  } catch (error) {
    wsSuccess.value = false
    wsResult.value = `✗ エラー: ${error.message}`
    wsLoading.value = false
  }
}
</script>
