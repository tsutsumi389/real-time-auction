/**
 * WebSocket Service
 * オークションライブ画面用のWebSocket接続管理
 */

class WebSocketService {
  constructor() {
    this.ws = null
    this.url = import.meta.env.VITE_WS_URL || 'ws://localhost/ws'
    this.token = null
    this.auctionId = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 3000 // 3秒
    this.pingInterval = 30000 // 30秒
    this.pingTimer = null
    this.isIntentionalClose = false

    // イベントハンドラー
    this.eventHandlers = {
      'auction:started': [],
      'auction:ended': [],
      'auction:cancelled': [],
      'item:started': [],
      'item:ended': [],
      'price:opened': [],
      'bid:placed': [],
      'participant:joined': [],
      'participant:left': [],
      'error': [],
      'connected': [],
      'disconnected': [],
      'reconnecting': [],
    }
  }

  /**
   * WebSocket接続を確立
   * @param {string} token - JWT認証トークン
   * @param {string} auctionId - オークションID
   */
  connect(token, auctionId) {
    this.token = token
    this.auctionId = auctionId
    this.isIntentionalClose = false

    const wsUrl = `${this.url}?token=${encodeURIComponent(token)}&auction_id=${encodeURIComponent(auctionId)}`

    try {
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        console.log('[WebSocket] Connected to auction:', auctionId)
        this.reconnectAttempts = 0
        this.startPingTimer()
        this.emit('connected', { auctionId })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.handleMessage(data)
        } catch (error) {
          console.error('[WebSocket] Failed to parse message:', error)
        }
      }

      this.ws.onerror = (error) => {
        console.error('[WebSocket] Error:', error)
        this.emit('error', { message: 'WebSocket通信エラーが発生しました' })
      }

      this.ws.onclose = (event) => {
        console.log('[WebSocket] Connection closed:', event.code, event.reason)
        this.stopPingTimer()
        this.emit('disconnected', { code: event.code, reason: event.reason })

        // 意図的なクローズでない場合は自動再接続を試みる
        if (!this.isIntentionalClose && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++
          console.log(`[WebSocket] Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
          this.emit('reconnecting', { attempt: this.reconnectAttempts, max: this.maxReconnectAttempts })

          setTimeout(() => {
            this.connect(this.token, this.auctionId)
          }, this.reconnectDelay)
        } else if (this.reconnectAttempts >= this.maxReconnectAttempts) {
          this.emit('error', { message: '接続の再試行回数が上限に達しました。ページを再読み込みしてください。' })
        }
      }
    } catch (error) {
      console.error('[WebSocket] Failed to connect:', error)
      this.emit('error', { message: 'WebSocket接続に失敗しました' })
    }
  }

  /**
   * WebSocket接続を切断
   */
  disconnect() {
    this.isIntentionalClose = true
    this.stopPingTimer()

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.close(1000, 'Client disconnect')
    }

    this.ws = null
    this.token = null
    this.auctionId = null
    this.reconnectAttempts = 0
  }

  /**
   * メッセージを処理
   * @param {object} data - WebSocketメッセージデータ
   */
  handleMessage(data) {
    const { type, payload } = data

    if (!type) {
      console.warn('[WebSocket] Received message without type:', data)
      return
    }

    console.log('[WebSocket] Received:', type, payload)

    // Ping/Pong処理
    if (type === 'ping') {
      this.send({ type: 'pong' })
      return
    }

    if (type === 'pong') {
      return
    }

    // イベントを発火
    this.emit(type, payload)
  }

  /**
   * サーバーにメッセージを送信
   * @param {object} data - 送信するデータ
   */
  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify(data))
        console.log('[WebSocket] Sent:', data.type)
      } catch (error) {
        console.error('[WebSocket] Failed to send message:', error)
        this.emit('error', { message: 'メッセージの送信に失敗しました' })
      }
    } else {
      console.warn('[WebSocket] Cannot send message, connection is not open')
      this.emit('error', { message: 'WebSocketが接続されていません' })
    }
  }

  /**
   * Pingタイマーを開始
   */
  startPingTimer() {
    this.stopPingTimer()
    this.pingTimer = setInterval(() => {
      this.send({ type: 'ping' })
    }, this.pingInterval)
  }

  /**
   * Pingタイマーを停止
   */
  stopPingTimer() {
    if (this.pingTimer) {
      clearInterval(this.pingTimer)
      this.pingTimer = null
    }
  }

  /**
   * イベントハンドラーを登録
   * @param {string} eventType - イベントタイプ
   * @param {Function} handler - ハンドラー関数
   */
  on(eventType, handler) {
    if (!this.eventHandlers[eventType]) {
      this.eventHandlers[eventType] = []
    }
    this.eventHandlers[eventType].push(handler)
  }

  /**
   * イベントハンドラーを解除
   * @param {string} eventType - イベントタイプ
   * @param {Function} handler - ハンドラー関数
   */
  off(eventType, handler) {
    if (!this.eventHandlers[eventType]) {
      return
    }
    this.eventHandlers[eventType] = this.eventHandlers[eventType].filter(h => h !== handler)
  }

  /**
   * イベントを発火
   * @param {string} eventType - イベントタイプ
   * @param {object} payload - イベントペイロード
   */
  emit(eventType, payload) {
    if (!this.eventHandlers[eventType]) {
      return
    }
    this.eventHandlers[eventType].forEach(handler => {
      try {
        handler(payload)
      } catch (error) {
        console.error(`[WebSocket] Error in ${eventType} handler:`, error)
      }
    })
  }

  /**
   * 接続状態を取得
   * @returns {boolean} 接続中の場合true
   */
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }

  /**
   * 接続状態を取得（詳細）
   * @returns {string} 接続状態 (connecting/open/closing/closed)
   */
  getReadyState() {
    if (!this.ws) {
      return 'closed'
    }

    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting'
      case WebSocket.OPEN:
        return 'open'
      case WebSocket.CLOSING:
        return 'closing'
      case WebSocket.CLOSED:
        return 'closed'
      default:
        return 'unknown'
    }
  }
}

// シングルトンインスタンスをエクスポート
export default new WebSocketService()
