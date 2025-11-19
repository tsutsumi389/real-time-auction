<template>
  <div class="home">
    <div class="welcome-card">
      <h2>🎉 開発環境セットアップ完了</h2>
      <p>Docker + Vite + Vue.js 3 で動作中</p>
    </div>

    <div class="status-grid">
      <div class="status-card">
        <h3>📡 API接続テスト</h3>
        <button @click="testApi" :disabled="apiLoading">
          {{ apiLoading ? 'テスト中...' : 'APIをテスト' }}
        </button>
        <div v-if="apiResult" class="result" :class="{ success: apiSuccess, error: !apiSuccess }">
          {{ apiResult }}
        </div>
      </div>

      <div class="status-card">
        <h3>🔌 WebSocket接続テスト</h3>
        <button @click="testWebSocket" :disabled="wsLoading">
          {{ wsLoading ? 'テスト中...' : 'WebSocketをテスト' }}
        </button>
        <div v-if="wsResult" class="result" :class="{ success: wsSuccess, error: !wsSuccess }">
          {{ wsResult }}
        </div>
      </div>
    </div>

    <div class="info-section">
      <h3>📋 次のステップ</h3>
      <ul>
        <li>データベーススキーマの設計と実装</li>
        <li>JWT認証システムの実装</li>
        <li>ユーザー管理APIの実装</li>
        <li>オークション管理機能の実装</li>
        <li>WebSocketリアルタイム通信の実装</li>
        <li>フロントエンドUIコンポーネントの構築</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const apiLoading = ref(false)
const apiResult = ref('')
const apiSuccess = ref(false)

const wsLoading = ref(false)
const wsResult = ref('')
const wsSuccess = ref(false)

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost/api'
const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost/ws'

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
    
    ws.onerror = (error) => {
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

<style scoped>
.home {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.welcome-card h2 {
  margin: 0 0 0.5rem;
  color: #667eea;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.status-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.status-card h3 {
  margin: 0 0 1rem;
  color: #333;
}

button {
  width: 100%;
  padding: 0.75rem;
  font-size: 1rem;
  border: none;
  border-radius: 4px;
  background: #667eea;
  color: white;
  cursor: pointer;
  transition: background 0.3s;
}

button:hover:not(:disabled) {
  background: #5568d3;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.result {
  margin-top: 1rem;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.result.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.result.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.info-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.info-section h3 {
  margin: 0 0 1rem;
  color: #333;
}

.info-section ul {
  margin: 0;
  padding-left: 1.5rem;
}

.info-section li {
  margin-bottom: 0.5rem;
  color: #666;
}
</style>
