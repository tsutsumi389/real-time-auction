/**
 * Bidder Auction Live Store
 * 入札者向けオークションライブ画面の状態管理とWebSocketイベント処理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getAuctionDetail,
  getPoints,
  placeBid as apiPlaceBid,
  getBidHistory as apiGetBidHistory,
} from '@/services/bidderBidApi'
import websocketService from '@/services/websocketService'
import { getToken, getUserFromToken } from '@/services/token'
import { useToast } from '@/composables/useToast'

export const useBidderAuctionLiveStore = defineStore('bidderAuctionLive', () => {
  // Toast composable
  const toast = useToast()

  // State
  const auction = ref(null)
  const items = ref([])
  const currentItem = ref(null)
  const bids = ref([])
  const points = ref({
    total: 0,
    available: 0,
    reserved: 0,
  })
  const loading = ref(false)
  const loadingBid = ref(false)
  const error = ref(null)
  const wsConnected = ref(false)
  const wsReconnecting = ref(false)
  const reconnectAttempt = ref(0)
  const maxReconnectAttempts = ref(5)
  const hasBidAtCurrentPrice = ref(false) // 現在価格で既に入札があるか

  // Computed

  /**
   * 現在の商品の価格
   */
  const currentPrice = computed(() => {
    return currentItem.value?.current_price || 0
  })

  /**
   * 現在の価格で入札可能なポイントがあるか
   */
  const hasEnoughPoints = computed(() => {
    return points.value.available >= currentPrice.value
  })

  /**
   * 現在のユーザーが勝者入札者か
   */
  const isOwnBidWinning = computed(() => {
    if (!currentItem.value || !bids.value.length) {
      return false
    }

    // 現在のユーザーIDを取得
    const user = getUserFromToken('bidder')
    if (!user || !user.bidderId) {
      return false
    }

    // is_winningがtrueの入札を探し、それが自分の入札かを確認
    const winningBid = bids.value.find((bid) => bid.is_winning === true)
    if (!winningBid) {
      return false
    }

    // bidder_idを比較（UUIDの文字列比較）
    return winningBid.bidder_id === user.bidderId
  })

  /**
   * 入札可能か（条件: 商品が開始済み、終了していない、ポイント十分、すでに勝者でない、現在価格で入札がない）
   */
  const canBid = computed(() => {
    if (!currentItem.value) {
      return false
    }
    const isStarted = currentItem.value.status === 'active'
    const hasPrice = currentPrice.value > 0
    return (
      isStarted &&
      hasPrice &&
      hasEnoughPoints.value &&
      !isOwnBidWinning.value &&
      !hasBidAtCurrentPrice.value &&
      !loadingBid.value
    )
  })

  /**
   * 入札不可能な理由
   */
  const cannotBidReason = computed(() => {
    if (!currentItem.value) {
      return '商品が選択されていません'
    }
    if (currentItem.value.status === 'pending') {
      return '商品はまだ開始されていません'
    }
    if (currentItem.value.status === 'ended') {
      return '商品は終了しました'
    }
    if (currentPrice.value === 0) {
      return '価格が開示されていません'
    }
    if (isOwnBidWinning.value) {
      return '現在あなたが最高入札者です'
    }
    if (hasBidAtCurrentPrice.value) {
      return '現在の価格で既に入札があります。次の価格開示をお待ちください'
    }
    if (!hasEnoughPoints.value) {
      return 'ポイントが不足しています'
    }
    if (loadingBid.value) {
      return '入札処理中...'
    }
    return ''
  })

  // Actions

  /**
   * 初期化: オークション詳細とポイント取得
   * @param {string} auctionId - オークションID
   */
  async function initialize(auctionId) {
    loading.value = true
    error.value = null

    try {
      // オークション詳細とポイントを並列取得
      const [auctionData, pointsData] = await Promise.all([
        getAuctionDetail(auctionId),
        getPoints(),
      ])

      // APIレスポンスはフラット構造（{ id, title, ..., items }）で返される
      // items を抽出してから auction に設定
      const { items: auctionItems, ...auctionInfo } = auctionData
      auction.value = auctionInfo
      items.value = auctionItems || []

      // ポイントAPIもフラット構造（{ total_points, available_points, reserved_points }）
      points.value = {
        total: pointsData.total_points || 0,
        available: pointsData.available_points || 0,
        reserved: pointsData.reserved_points || 0,
      }

      // 最初の商品を選択（status=activeの商品があればそれを、なければ最初の商品）
      const activeItem = items.value.find((item) => item.status === 'active')
      if (activeItem) {
        await selectItem(activeItem.id)
      } else if (items.value.length > 0) {
        currentItem.value = items.value[0]
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || 'オークション情報の取得に失敗しました'
      console.error('[bidderAuctionLive] Initialize error:', err)
      return false
    }
  }

  /**
   * 商品を選択
   * @param {number} itemId - 商品ID
   */
  async function selectItem(itemId) {
    const item = items.value.find((i) => i.id === itemId)
    if (!item) {
      console.warn('[bidderAuctionLive] Item not found:', itemId)
      return
    }

    currentItem.value = item
    bids.value = []
    hasBidAtCurrentPrice.value = false // 商品切り替え時にリセット

    // 入札履歴を取得
    await fetchBidHistory(itemId)
  }

  /**
   * 入札履歴を取得
   * @param {number} itemId - 商品ID
   */
  async function fetchBidHistory(itemId) {
    try {
      const response = await apiGetBidHistory(itemId, { limit: 50, offset: 0 })
      bids.value = response.bids || []

      // 現在価格で既に入札があるかチェック
      if (currentItem.value && currentItem.value.current_price) {
        const currentPrice = currentItem.value.current_price
        const hasBidAtPrice = bids.value.some((bid) => bid.price === currentPrice)
        hasBidAtCurrentPrice.value = hasBidAtPrice
      }
    } catch (err) {
      console.error('[bidderAuctionLive] Fetch bid history error:', err)
      // エラーは表示せず、履歴が空になるだけ
    }
  }

  /**
   * 入札実行
   * @param {number} itemId - 商品ID
   * @param {number} price - 入札価格
   */
  async function placeBid(itemId, price) {
    loadingBid.value = true
    error.value = null

    try {
      const response = await apiPlaceBid(itemId, price)

      // ポイントを更新（APIレスポンス構造を正規化）
      if (response.points) {
        points.value = {
          total: response.points.total_points ?? response.points.total ?? 0,
          available: response.points.available_points ?? response.points.available ?? 0,
          reserved: response.points.reserved_points ?? response.points.reserved ?? 0,
        }
      }

      // 楽観的UI更新: 入札を履歴の先頭に追加（WebSocketイベントで重複追加されないよう、IDをチェック）
      if (response.bid) {
        const existingBid = bids.value.find((b) => b.id === response.bid.id)
        if (!existingBid) {
          bids.value.unshift(response.bid)
        }
      }

      loadingBid.value = false
      return { success: true, data: response }
    } catch (err) {
      loadingBid.value = false

      // エラーメッセージを整形
      let errorMessage = 'エラーが発生しました'
      if (err.message) {
        errorMessage = err.message
      } else if (err.response?.data?.error) {
        errorMessage = err.response.data.error
      }

      error.value = errorMessage
      console.error('[bidderAuctionLive] Place bid error:', err)
      return { success: false, error: errorMessage }
    }
  }

  /**
   * WebSocket接続
   * @param {string} auctionId - オークションID
   */
  function connectWebSocket(auctionId) {
    const token = getToken('bidder')
    if (!token) {
      console.error('[bidderAuctionLive] Cannot connect WebSocket: No token')
      error.value = '認証トークンがありません'
      return
    }

    // イベントハンドラーを登録
    websocketService.on('connected', onWebSocketConnected)
    websocketService.on('disconnected', onWebSocketDisconnected)
    websocketService.on('reconnecting', onWebSocketReconnecting)
    websocketService.on('error', onWebSocketError)
    websocketService.on('price:opened', onPriceOpened)
    websocketService.on('bid:placed', onBidPlaced)
    websocketService.on('item:started', onItemStarted)
    websocketService.on('item:ended', onItemEnded)
    websocketService.on('auction:ended', onAuctionEnded)

    // 接続
    websocketService.connect(token, auctionId)
  }

  /**
   * WebSocket切断
   */
  function disconnectWebSocket() {
    // イベントハンドラーを解除
    websocketService.off('connected', onWebSocketConnected)
    websocketService.off('disconnected', onWebSocketDisconnected)
    websocketService.off('reconnecting', onWebSocketReconnecting)
    websocketService.off('error', onWebSocketError)
    websocketService.off('price:opened', onPriceOpened)
    websocketService.off('bid:placed', onBidPlaced)
    websocketService.off('item:started', onItemStarted)
    websocketService.off('item:ended', onItemEnded)
    websocketService.off('auction:ended', onAuctionEnded)

    // 切断
    websocketService.disconnect()
  }

  // WebSocketイベントハンドラー

  function onWebSocketConnected() {
    console.log('[bidderAuctionLive] WebSocket connected')
    wsConnected.value = true
    wsReconnecting.value = false
    reconnectAttempt.value = 0

    // オークションルームに参加
    if (auction.value?.id) {
      websocketService.send({
        type: 'subscribe',
        data: {
          auction_id: auction.value.id
        }
      })
      console.log('[bidderAuctionLive] Sent subscribe event for auction:', auction.value.id)
    }
  }

  function onWebSocketDisconnected() {
    console.log('[bidderAuctionLive] WebSocket disconnected')
    wsConnected.value = false
  }

  function onWebSocketReconnecting(payload) {
    console.log('[bidderAuctionLive] WebSocket reconnecting:', payload)
    wsReconnecting.value = true
    reconnectAttempt.value = payload.attempt
    maxReconnectAttempts.value = payload.max
  }

  function onWebSocketError(payload) {
    console.error('[bidderAuctionLive] WebSocket error:', payload)
    error.value = payload.message
  }

  /**
   * 価格開示イベント
   * @param {object} payload - { item_id, price, opened_at } or { payload: { item_id, price } }
   */
  function onPriceOpened(payload) {
    console.log('[bidderAuctionLive] Price opened:', payload)

    // payloadから必要なデータを取得（ネストされた構造にも対応）
    const item_id = payload.item_id || payload.payload?.item_id
    const price = payload.price || payload.payload?.price
    const opened_at = payload.opened_at || payload.payload?.price_history?.disclosed_at

    if (!item_id || !price) {
      console.error('[bidderAuctionLive] Invalid price:opened payload:', payload)
      return
    }

    // 商品の価格を更新
    const item = items.value.find((i) => i.id === item_id)
    if (item) {
      item.current_price = price
      if (opened_at) {
        item.updated_at = opened_at
      }
    }

    // 現在の商品なら更新し、通知表示
    if (currentItem.value?.id === item_id) {
      currentItem.value.current_price = price
      if (opened_at) {
        currentItem.value.updated_at = opened_at
      }

      // 新しい価格が開示されたので、入札可能状態にリセット
      hasBidAtCurrentPrice.value = false

      // トースト通知
      toast.info(
        '価格が開示されました',
        `新しい価格: ${new Intl.NumberFormat('ja-JP').format(price)}ポイント`,
        3000
      )
    }
  }

  /**
   * 入札イベント
   * @param {object} payload - { item_id, bid, points }
   */
  function onBidPlaced(payload) {
    console.log('[bidderAuctionLive] Bid placed:', payload)

    // 自分の入札の場合はポイントを更新
    if (payload.points) {
      // APIレスポンス構造を正規化（total_points -> total など）
      points.value = {
        total: payload.points.total_points ?? payload.points.total ?? 0,
        available: payload.points.available_points ?? payload.points.available ?? 0,
        reserved: payload.points.reserved_points ?? payload.points.reserved ?? 0,
      }
    }

    // 入札履歴に追加（重複チェック）
    if (payload.bid && currentItem.value?.id === payload.item_id) {
      const existingBid = bids.value.find((b) => b.id === payload.bid.id)
      if (!existingBid) {
        bids.value.unshift(payload.bid)

        // 現在価格で入札があったことを記録
        hasBidAtCurrentPrice.value = true

        // 他者の入札通知（自分の入札はhandlePlaceBidで通知済み）
        if (!payload.points && payload.bid.bidder_display_name) {
          toast.warning(
            '他の入札者が入札しました',
            `${payload.bid.bidder_display_name}さんが ${new Intl.NumberFormat('ja-JP').format(payload.bid.price)}ポイントで入札`,
            3000
          )
        }
      }
    }
  }

  /**
   * 商品開始イベント
   * @param {object} payload - { item: { id, auction_id, name, current_price, started_at, status } }
   */
  function onItemStarted(payload) {
    console.log('[bidderAuctionLive] Item started:', payload)

    // payloadから商品情報を取得
    const itemData = payload.item
    if (!itemData || !itemData.id) {
      console.error('[bidderAuctionLive] Invalid item:started payload:', payload)
      return
    }

    // 商品リストを更新
    const item = items.value.find((i) => i.id === itemData.id)
    if (item) {
      item.status = itemData.status || 'active'
      item.current_price = itemData.current_price
      if (itemData.started_at) {
        item.started_at = itemData.started_at
      }

      // 通知
      toast.info(
        '商品が開始されました',
        `${item.name}の入札が開始されました`,
        3000
      )
    }

    // 現在の商品を更新
    if (currentItem.value?.id === itemData.id) {
      currentItem.value.status = itemData.status || 'active'
      currentItem.value.current_price = itemData.current_price
      if (itemData.started_at) {
        currentItem.value.started_at = itemData.started_at
      }
    }
  }

  /**
   * 商品終了イベント
   * @param {object} payload - { item_id, status, winner_id, ended_at, points } or { payload: { item } }
   */
  function onItemEnded(payload) {
    console.log('[bidderAuctionLive] Item ended:', payload)

    // payloadから必要なデータを取得（ネストされた構造にも対応）
    const item_id = payload.item_id || payload.payload?.item?.id
    const status = payload.status || payload.payload?.item?.status || 'ended'
    const winner_id = payload.winner_id || payload.payload?.item?.winner_id
    const ended_at = payload.ended_at || payload.payload?.item?.ended_at

    if (!item_id) {
      console.error('[bidderAuctionLive] Invalid item:ended payload:', payload)
      return
    }

    const item = items.value.find((i) => i.id === item_id)
    if (item) {
      item.status = status
      if (winner_id !== undefined) {
        item.winner_id = winner_id
      }
      if (ended_at) {
        item.ended_at = ended_at
      }
    }

    if (currentItem.value?.id === item_id) {
      currentItem.value.status = status
      if (winner_id !== undefined) {
        currentItem.value.winner_id = winner_id
      }
      if (ended_at) {
        currentItem.value.ended_at = ended_at
      }
    }

    // 落札/非落札によるポイント更新と通知
    if (payload.points) {
      const oldPoints = points.value
      // APIレスポンス構造を正規化
      const normalizedPoints = {
        total: payload.points.total_points ?? payload.points.total ?? 0,
        available: payload.points.available_points ?? payload.points.available ?? 0,
        reserved: payload.points.reserved_points ?? payload.points.reserved ?? 0,
      }
      points.value = normalizedPoints

      // ポイントが消費された場合（落札）
      if (normalizedPoints.reserved < oldPoints.reserved) {
        const targetItem = items.value.find((i) => i.id === item_id)
        toast.success(
          '落札しました！',
          `${targetItem?.name || '商品'}を落札しました`,
          5000
        )
      } else if (normalizedPoints.available > oldPoints.available) {
        // ポイントが戻った場合（非落札）
        toast.info(
          '商品が終了しました',
          '予約ポイントが返却されました',
          3000
        )
      }
    }
  }

  /**
   * オークション終了イベント
   * @param {object} payload - { auction_id, status, ended_at }
   */
  function onAuctionEnded(payload) {
    console.log('[bidderAuctionLive] Auction ended:', payload)

    if (auction.value && auction.value.id === payload.auction_id) {
      auction.value.status = payload.status
      auction.value.ended_at = payload.ended_at

      // 通知
      toast.info(
        'オークションが終了しました',
        `${auction.value.title}が終了しました。ご参加ありがとうございました。`,
        5000
      )
    }
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  /**
   * ストアをリセット
   */
  function reset() {
    auction.value = null
    items.value = []
    currentItem.value = null
    bids.value = []
    points.value = { total: 0, available: 0, reserved: 0 }
    loading.value = false
    loadingBid.value = false
    error.value = null
    wsConnected.value = false
    wsReconnecting.value = false
    reconnectAttempt.value = 0
    maxReconnectAttempts.value = 5
    hasBidAtCurrentPrice.value = false
  }

  return {
    // State
    auction,
    items,
    currentItem,
    bids,
    points,
    loading,
    loadingBid,
    error,
    wsConnected,
    wsReconnecting,
    reconnectAttempt,
    maxReconnectAttempts,
    // Computed
    currentPrice,
    hasEnoughPoints,
    isOwnBidWinning,
    canBid,
    cannotBidReason,
    // Actions
    initialize,
    selectItem,
    fetchBidHistory,
    placeBid,
    connectWebSocket,
    disconnectWebSocket,
    clearError,
    reset,
  }
})
