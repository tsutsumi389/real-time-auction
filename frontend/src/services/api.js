/**
 * API Client
 * Axiosインスタンスの設定とインターセプターの実装
 */
import axios from 'axios'
import { getToken, removeToken } from './token'

// 環境変数からAPIベースURLを取得
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost/api'

// Axiosインスタンスの作成
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// リクエストインターセプター
apiClient.interceptors.request.use(
  (config) => {
    // トークンが存在する場合、Authorizationヘッダーに追加
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// レスポンスインターセプター
apiClient.interceptors.response.use(
  (response) => {
    // 成功レスポンスはそのまま返す
    return response
  },
  (error) => {
    // エラーハンドリング
    if (error.response) {
      // サーバーからのエラーレスポンス
      const { status, data } = error.response

      switch (status) {
        case 401:
          // 認証エラー: トークンを削除してログイン画面へ
          removeToken()
          // ログイン画面へのリダイレクトはルーターで処理
          break
        case 403:
          // 権限エラー
          console.error('Access forbidden:', data.error)
          break
        case 404:
          // リソースが見つからない
          console.error('Resource not found:', data.error)
          break
        case 500:
          // サーバーエラー
          console.error('Server error:', data.error)
          break
        default:
          console.error('API error:', data.error)
      }

      // エラーメッセージを統一形式で返す
      return Promise.reject({
        status,
        message: data.error || 'An error occurred',
        details: data.details || null,
      })
    } else if (error.request) {
      // リクエストは送信されたがレスポンスがない（ネットワークエラー）
      console.error('Network error:', error.message)
      return Promise.reject({
        status: 0,
        message: 'Network error. Please check your connection.',
        details: null,
      })
    } else {
      // リクエスト設定時のエラー
      console.error('Request error:', error.message)
      return Promise.reject({
        status: -1,
        message: error.message,
        details: null,
      })
    }
  }
)

export default apiClient
