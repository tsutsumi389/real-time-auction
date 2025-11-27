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
import { getToken } from '@/services/token'
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
    // 最新の入札が自分のものか確認
    const latestBid = bids.value[0]
    return latestBid?.is_winning === true
  })

  /**
   * 入札可能か（条件: 商品が開始済み、終了していない、ポイント十分、すでに勝者でない）
   */
  const canBid = computed(() => {
    if (!currentItem.value) {
      return false
    }
    const isStarted = currentItem.value.status === 'started'
    const hasPrice = currentPrice.value > 0
    return (
      isStarted &&
      hasPrice &&
      hasEnoughPoints.value &&
      !isOwnBidWinning.value &&
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

      auction.value = auctionData.auction
      items.value = auctionData.items || []
      points.value = pointsData.points || { total: 0, available: 0, reserved: 0 }

      // 最初の商品を選択（status=startedの商品があればそれを、なければ最初の商品）
      const startedItem = items.value.find((item) => item.status === 'started')
      if (startedItem) {
        await selectItem(startedItem.id)
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

      // ポイントを更新
      if (response.points) {
        points.value = response.points
      }

      // 楽観的UI更新: 入札を履歴の先頭に追加
      if (response.bid) {
        bids.value.unshift(response.bid)
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
   * @param {object} payload - { item_id, price, opened_at }
   */
  function onPriceOpened(payload) {
    console.log('[bidderAuctionLive] Price opened:', payload)

    // 商品の価格を更新
    const item = items.value.find((i) => i.id === payload.item_id)
    if (item) {
      item.current_price = payload.price
      item.updated_at = payload.opened_at
    }

    // 現在の商品なら更新し、通知表示
    if (currentItem.value?.id === payload.item_id) {
      currentItem.value.current_price = payload.price
      currentItem.value.updated_at = payload.opened_at

      // トースト通知
      toast.info(
        '価格が開示されました',
        `新しい価格: ${new Intl.NumberFormat('ja-JP').format(payload.price)}ポイント`,
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
      points.value = payload.points
    }

    // 入札履歴に追加（重複チェック）
    if (payload.bid && currentItem.value?.id === payload.item_id) {
      const existingBid = bids.value.find((b) => b.id === payload.bid.id)
      if (!existingBid) {
        bids.value.unshift(payload.bid)

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
   * @param {object} payload - { item_id, status, started_at }
   */
  function onItemStarted(payload) {
    console.log('[bidderAuctionLive] Item started:', payload)

    const item = items.value.find((i) => i.id === payload.item_id)
    if (item) {
      item.status = payload.status
      item.started_at = payload.started_at

      // 通知
      toast.info(
        '商品が開始されました',
        `${item.name}の入札が開始されました`,
        3000
      )
    }

    if (currentItem.value?.id === payload.item_id) {
      currentItem.value.status = payload.status
      currentItem.value.started_at = payload.started_at
    }
  }

  /**
   * 商品終了イベント
   * @param {object} payload - { item_id, status, winner_id, ended_at, points }
   */
  function onItemEnded(payload) {
    console.log('[bidderAuctionLive] Item ended:', payload)

    const item = items.value.find((i) => i.id === payload.item_id)
    if (item) {
      item.status = payload.status
      item.winner_id = payload.winner_id
      item.ended_at = payload.ended_at
    }

    if (currentItem.value?.id === payload.item_id) {
      currentItem.value.status = payload.status
      currentItem.value.winner_id = payload.winner_id
      currentItem.value.ended_at = payload.ended_at
    }

    // 落札/非落札によるポイント更新と通知
    if (payload.points) {
      const oldPoints = points.value
      points.value = payload.points

      // ポイントが消費された場合（落札）
      if (payload.points.reserved < oldPoints.reserved) {
        const item = items.value.find((i) => i.id === payload.item_id)
        toast.success(
          '落札しました！',
          `${item?.name || '商品'}を落札しました`,
          5000
        )
      } else if (payload.points.available > oldPoints.available) {
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
