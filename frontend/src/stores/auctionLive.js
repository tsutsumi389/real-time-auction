/**
 * Auction Live Store
 * オークションライブ画面の状態管理（主催者用）
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getAuctionDetail,
  getParticipants,
  startItem,
  openPrice,
  endItem,
  getBidHistory,
  getPriceHistory,
  endAuction,
  cancelAuction,
} from '@/services/auctionApi'
import websocketService from '@/services/websocketService'

export const useAuctionLiveStore = defineStore('auctionLive', () => {
  // State
  const auction = ref(null)
  const items = ref([])
  const currentItem = ref(null)
  const participants = ref([])
  const bids = ref([])
  const priceHistory = ref([])
  const loading = ref(false)
  const error = ref(null)
  const wsConnected = ref(false)
  const wsReconnecting = ref(false)

  // Computed
  const activeItem = computed(() => {
    return items.value.find(item => item.status === 'active')
  })

  const pendingItems = computed(() => {
    return items.value.filter(item => item.status === 'pending')
  })

  const endedItems = computed(() => {
    return items.value.filter(item => item.status === 'ended')
  })

  const currentPrice = computed(() => {
    return currentItem.value?.current_price || 0
  })

  const participantCount = computed(() => {
    return participants.value.length
  })

  const activeParticipants = computed(() => {
    return participants.value.filter(p => p.status === 'active')
  })

  const latestBid = computed(() => {
    return bids.value.length > 0 ? bids.value[0] : null
  })

  // Actions

  /**
   * オークション詳細を取得して初期化
   * @param {string} auctionId - オークションID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function initialize(auctionId) {
    loading.value = true
    error.value = null

    try {
      const response = await getAuctionDetail(auctionId)

      // APIレスポンスは { id, title, ..., items: [] } の構造
      // itemsを取り出して、残りをauctionとして設定
      const { items: responseItems, ...auctionData } = response
      auction.value = auctionData
      items.value = responseItems || []

      // 最初のアクティブまたはペンディングアイテムを選択
      const active = items.value.find(item => item.status === 'active')
      const pending = items.value.find(item => item.status === 'pending')
      currentItem.value = active || pending || items.value[0] || null

      // 参加者一覧を取得
      await fetchParticipants(auctionId)

      // 現在のアイテムの入札履歴と価格履歴を取得
      if (currentItem.value) {
        await Promise.all([
          fetchBidHistory(currentItem.value.id),
          fetchPriceHistory(currentItem.value.id),
        ])
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || 'オークション情報の取得に失敗しました'
      return false
    }
  }

  /**
   * 参加者一覧を取得
   * @param {string} auctionId - オークションID
   */
  async function fetchParticipants(auctionId) {
    try {
      const response = await getParticipants(auctionId)
      participants.value = response.participants || []
    } catch (err) {
      console.error('Failed to fetch participants:', err)
    }
  }

  /**
   * 入札履歴を取得
   * @param {string} itemId - 商品ID
   */
  async function fetchBidHistory(itemId) {
    try {
      const response = await getBidHistory(itemId)
      bids.value = response.bids || []
    } catch (err) {
      console.error('Failed to fetch bid history:', err)
    }
  }

  /**
   * 価格開示履歴を取得
   * @param {string} itemId - 商品ID
   */
  async function fetchPriceHistory(itemId) {
    try {
      const response = await getPriceHistory(itemId)
      priceHistory.value = response.price_history || []
    } catch (err) {
      console.error('Failed to fetch price history:', err)
    }
  }

  /**
   * 商品を開始
   * @param {string} itemId - 商品ID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleStartItem(itemId) {
    error.value = null

    try {
      const response = await startItem(itemId)

      // アイテムのステータスを更新
      const index = items.value.findIndex(item => item.id === itemId)
      if (index !== -1) {
        items.value[index] = { ...items.value[index], ...response.item }
        currentItem.value = items.value[index]
      }

      // 入札履歴と価格履歴をリセット
      bids.value = []
      priceHistory.value = []
      await fetchPriceHistory(itemId)

      return true
    } catch (err) {
      error.value = err.response?.data?.error || '商品の開始に失敗しました'
      return false
    }
  }

  /**
   * 価格を開示
   * @param {string} itemId - 商品ID
   * @param {number} price - 開示する価格
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleOpenPrice(itemId, price) {
    error.value = null

    try {
      const response = await openPrice(itemId, { new_price: price })

      // 価格履歴に追加
      priceHistory.value.unshift(response.price_history)

      // アイテムの現在価格を更新
      const index = items.value.findIndex(item => item.id === itemId)
      if (index !== -1) {
        items.value[index].current_price = price
        if (currentItem.value?.id === itemId) {
          currentItem.value.current_price = price
        }
      }

      return true
    } catch (err) {
      error.value = err.response?.data?.error || '価格の開示に失敗しました'
      return false
    }
  }

  /**
   * 商品を終了
   * @param {string} itemId - 商品ID
   * @returns {Promise<object|null>} 成功した場合は落札情報、失敗した場合はnull
   */
  async function handleEndItem(itemId) {
    error.value = null

    try {
      const response = await endItem(itemId)

      // アイテムのステータスを更新
      const index = items.value.findIndex(item => item.id === itemId)
      if (index !== -1) {
        items.value[index] = { ...items.value[index], ...response.item }
      }

      // 次のペンディングアイテムを選択
      const nextPending = items.value.find(item => item.status === 'pending')
      if (nextPending) {
        currentItem.value = nextPending
        await Promise.all([
          fetchBidHistory(nextPending.id),
          fetchPriceHistory(nextPending.id),
        ])
      } else {
        currentItem.value = null
        bids.value = []
        priceHistory.value = []
      }

      return response.winner || null
    } catch (err) {
      error.value = err.response?.data?.error || '商品の終了に失敗しました'
      return null
    }
  }

  /**
   * オークションを終了
   * @param {string} auctionId - オークションID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleEndAuction(auctionId) {
    error.value = null

    try {
      const response = await endAuction(auctionId)
      auction.value = { ...auction.value, ...response.auction }
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'オークションの終了に失敗しました'
      return false
    }
  }

  /**
   * オークションを中止（system_adminのみ）
   * @param {string} auctionId - オークションID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleCancelAuction(auctionId) {
    error.value = null

    try {
      const response = await cancelAuction(auctionId)
      auction.value = { ...auction.value, ...response.auction }
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'オークションの中止に失敗しました'
      return false
    }
  }

  /**
   * 商品を選択
   * @param {string} itemId - 商品ID
   */
  async function selectItem(itemId) {
    const item = items.value.find(i => i.id === itemId)
    if (item) {
      currentItem.value = item
      await Promise.all([
        fetchBidHistory(itemId),
        fetchPriceHistory(itemId),
      ])
    }
  }

  /**
   * WebSocketイベント: 入札が行われた
   * @param {object} payload - イベントペイロード
   */
  function onBidPlaced(payload) {
    const { bid, item_id } = payload

    // 入札履歴に追加（先頭に挿入）
    bids.value.unshift(bid)

    // アイテムの最終入札価格を更新
    const index = items.value.findIndex(item => item.id === item_id)
    if (index !== -1) {
      items.value[index].last_bid_price = bid.price
      if (currentItem.value?.id === item_id) {
        currentItem.value.last_bid_price = bid.price
      }
    }
  }

  /**
   * WebSocketイベント: 価格が開示された
   * @param {object} payload - イベントペイロード
   */
  function onPriceOpened(payload) {
    const { price_history, item_id, price } = payload

    // 価格履歴に追加
    priceHistory.value.unshift(price_history)

    // アイテムの現在価格を更新
    const index = items.value.findIndex(item => item.id === item_id)
    if (index !== -1) {
      items.value[index].current_price = price
      if (currentItem.value?.id === item_id) {
        currentItem.value.current_price = price
      }
    }
  }

  /**
   * WebSocketイベント: 商品が開始された
   * @param {object} payload - イベントペイロード
   */
  function onItemStarted(payload) {
    const { item } = payload

    // アイテムのステータスを更新
    const index = items.value.findIndex(i => i.id === item.id)
    if (index !== -1) {
      items.value[index] = { ...items.value[index], ...item }
    }
  }

  /**
   * WebSocketイベント: 商品が終了された
   * @param {object} payload - イベントペイロード
   */
  function onItemEnded(payload) {
    const { item } = payload

    // アイテムのステータスを更新
    const index = items.value.findIndex(i => i.id === item.id)
    if (index !== -1) {
      items.value[index] = { ...items.value[index], ...item }
    }
  }

  /**
   * WebSocketイベント: 参加者が参加した
   * @param {object} payload - イベントペイロード
   */
  function onParticipantJoined(payload) {
    const { participant } = payload

    // 既に存在しない場合のみ追加
    const exists = participants.value.some(p => p.id === participant.id)
    if (!exists) {
      participants.value.push(participant)
    }
  }

  /**
   * WebSocketイベント: 参加者が退出した
   * @param {object} payload - イベントペイロード
   */
  function onParticipantLeft(payload) {
    const { bidder_id } = payload

    // 参加者のステータスを更新または削除
    const index = participants.value.findIndex(p => p.id === bidder_id)
    if (index !== -1) {
      participants.value.splice(index, 1)
    }
  }

  /**
   * WebSocketイベント: オークションが終了された
   * @param {object} payload - イベントペイロード
   */
  function onAuctionEnded(payload) {
    const { auction: updatedAuction } = payload
    auction.value = { ...auction.value, ...updatedAuction }
  }

  /**
   * WebSocketイベント: オークションが中止された
   * @param {object} payload - イベントペイロード
   */
  function onAuctionCancelled(payload) {
    const { auction: updatedAuction } = payload
    auction.value = { ...auction.value, ...updatedAuction }
  }

  /**
   * WebSocket接続を確立
   * @param {string} token - JWT認証トークン
   * @param {string} auctionId - オークションID
   */
  function connectWebSocket(token, auctionId) {
    // イベントハンドラーを登録
    websocketService.on('bid:placed', onBidPlaced)
    websocketService.on('price:opened', onPriceOpened)
    websocketService.on('item:started', onItemStarted)
    websocketService.on('item:ended', onItemEnded)
    websocketService.on('participant:joined', onParticipantJoined)
    websocketService.on('participant:left', onParticipantLeft)
    websocketService.on('auction:ended', onAuctionEnded)
    websocketService.on('auction:cancelled', onAuctionCancelled)

    websocketService.on('connected', () => {
      wsConnected.value = true
      wsReconnecting.value = false
    })

    websocketService.on('disconnected', () => {
      wsConnected.value = false
    })

    websocketService.on('reconnecting', () => {
      wsReconnecting.value = true
    })

    websocketService.on('error', (payload) => {
      error.value = payload.message
    })

    // 接続
    websocketService.connect(token, auctionId)
  }

  /**
   * WebSocket接続を切断
   */
  function disconnectWebSocket() {
    // イベントハンドラーを解除
    websocketService.off('bid:placed', onBidPlaced)
    websocketService.off('price:opened', onPriceOpened)
    websocketService.off('item:started', onItemStarted)
    websocketService.off('item:ended', onItemEnded)
    websocketService.off('participant:joined', onParticipantJoined)
    websocketService.off('participant:left', onParticipantLeft)
    websocketService.off('auction:ended', onAuctionEnded)
    websocketService.off('auction:cancelled', onAuctionCancelled)

    // 接続を切断
    websocketService.disconnect()
    wsConnected.value = false
    wsReconnecting.value = false
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
    participants.value = []
    bids.value = []
    priceHistory.value = []
    loading.value = false
    error.value = null
    wsConnected.value = false
    wsReconnecting.value = false
  }

  return {
    // State
    auction,
    items,
    currentItem,
    participants,
    bids,
    priceHistory,
    loading,
    error,
    wsConnected,
    wsReconnecting,
    // Computed
    activeItem,
    pendingItems,
    endedItems,
    currentPrice,
    participantCount,
    activeParticipants,
    latestBid,
    // Actions
    initialize,
    fetchParticipants,
    fetchBidHistory,
    fetchPriceHistory,
    handleStartItem,
    handleOpenPrice,
    handleEndItem,
    handleEndAuction,
    handleCancelAuction,
    selectItem,
    connectWebSocket,
    disconnectWebSocket,
    clearError,
    reset,
  }
})
