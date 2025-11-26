<script setup>
import { onMounted, onUnmounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuctionLiveStore } from '@/stores/auctionLive'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import ItemInfo from '@/components/auction/ItemInfo.vue'
import ControlPanel from '@/components/auction/ControlPanel.vue'
import BidHistory from '@/components/auction/BidHistory.vue'
import ParticipantList from '@/components/auction/ParticipantList.vue'
import PriceHistoryList from '@/components/auction/PriceHistoryList.vue'
import ImageSlider from '@/components/auction/ImageSlider.vue'
import WinnerModal from '@/components/auction/WinnerModal.vue'
import Alert from '@/components/ui/Alert.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

const route = useRoute()
const router = useRouter()
const auctionLiveStore = useAuctionLiveStore()
const authStore = useAuthStore()
const toast = useToast()

const auctionId = computed(() => route.params.id)
const showWinnerModal = ref(false)
const winnerData = ref(null)

const wsStatus = computed(() => {
  if (auctionLiveStore.wsReconnecting) {
    return { text: '再接続中...', color: 'text-yellow-600' }
  }
  if (auctionLiveStore.wsConnected) {
    return { text: '接続中', color: 'text-green-600' }
  }
  return { text: '切断', color: 'text-red-600' }
})

const auctionStatusText = computed(() => {
  const status = auctionLiveStore.auction?.status
  switch (status) {
    case 'pending':
      return '待機中'
    case 'active':
      return '進行中'
    case 'ended':
      return '終了'
    case 'cancelled':
      return '中止'
    default:
      return status
  }
})

const auctionStatusClass = computed(() => {
  const status = auctionLiveStore.auction?.status
  switch (status) {
    case 'pending':
      return 'bg-gray-100 text-gray-700'
    case 'active':
      return 'bg-green-100 text-green-700'
    case 'ended':
      return 'bg-blue-100 text-blue-700'
    case 'cancelled':
      return 'bg-red-100 text-red-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
})

async function handleStartItem(itemId) {
  const success = await auctionLiveStore.handleStartItem(itemId)
  if (success) {
    toast.success('商品を開始しました')
  } else if (auctionLiveStore.error) {
    toast.error('商品開始エラー', auctionLiveStore.error)
  }
}

async function handleOpenPrice(itemId, price) {
  const success = await auctionLiveStore.handleOpenPrice(itemId, price)
  if (success) {
    toast.success('価格を開示しました', `${new Intl.NumberFormat('ja-JP').format(price)} pt`)
  } else if (auctionLiveStore.error) {
    toast.error('価格開示エラー', auctionLiveStore.error)
  }
}

async function handleEndItem(itemId) {
  const winner = await auctionLiveStore.handleEndItem(itemId)
  if (winner) {
    winnerData.value = {
      ...winner,
      item: auctionLiveStore.currentItem
    }
    showWinnerModal.value = true
    toast.success('商品が終了しました', `落札者: ${winner.display_name}`)
  } else if (auctionLiveStore.error) {
    toast.error('商品終了エラー', auctionLiveStore.error)
  }
}

async function handleEndAuction(auctionId) {
  const success = await auctionLiveStore.handleEndAuction(auctionId)
  if (success) {
    toast.success('オークションを終了しました', 'オークション一覧に戻ります')
    setTimeout(() => {
      router.push({ name: 'auction-list' })
    }, 2000)
  } else if (auctionLiveStore.error) {
    toast.error('オークション終了エラー', auctionLiveStore.error)
  }
}

async function handleCancelAuction(auctionId) {
  const success = await auctionLiveStore.handleCancelAuction(auctionId)
  if (success) {
    toast.warning('オークションを緊急停止しました', 'すべてのポイントが返金されました')
    setTimeout(() => {
      router.push({ name: 'auction-list' })
    }, 2000)
  } else if (auctionLiveStore.error) {
    toast.error('オークション中止エラー', auctionLiveStore.error)
  }
}

function handleSelectItem(itemId) {
  auctionLiveStore.selectItem(itemId)
}

onMounted(async () => {
  const success = await auctionLiveStore.initialize(auctionId.value)
  if (success) {
    const token = authStore.token
    if (token) {
      auctionLiveStore.connectWebSocket(token, auctionId.value)
      toast.info('オークションに接続しました', 'リアルタイム更新が有効です')
    }
  } else {
    toast.error('オークション初期化エラー', auctionLiveStore.error)
  }
})

onUnmounted(() => {
  auctionLiveStore.disconnectWebSocket()
  auctionLiveStore.reset()
})
</script>

<template>
  <div class="p-6 space-y-6">
    <!-- ヘッダー -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold">オークションライブ</h1>
        <div v-if="auctionLiveStore.auction" class="flex items-center gap-4 mt-2">
          <div class="text-lg text-gray-600">{{ auctionLiveStore.auction.title }}</div>
          <div :class="['px-3 py-1 rounded-full text-sm font-medium', auctionStatusClass]">
            {{ auctionStatusText }}
          </div>
        </div>
      </div>
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2">
          <div :class="['w-3 h-3 rounded-full', auctionLiveStore.wsConnected ? 'bg-green-500' : 'bg-red-500']"></div>
          <span :class="['text-sm font-medium', wsStatus.color]">{{ wsStatus.text }}</span>
        </div>
        <Button @click="router.push({ name: 'auction-list' })" variant="outline">
          一覧に戻る
        </Button>
      </div>
    </div>

    <!-- エラー表示 -->
    <Alert v-if="auctionLiveStore.error" variant="destructive">
      {{ auctionLiveStore.error }}
    </Alert>

    <!-- 落札者発表モーダル -->
    <WinnerModal
      :open="showWinnerModal"
      :winner="winnerData"
      :item="winnerData?.item"
      @close="showWinnerModal = false"
    />

    <!-- ローディング -->
    <div v-if="auctionLiveStore.loading" class="flex justify-center py-12">
      <LoadingSpinner />
    </div>

    <!-- メインコンテンツ -->
    <div v-else-if="auctionLiveStore.auction">
      <!-- モバイル・タブレット: 縦並び -->
      <div class="lg:hidden space-y-6">
        <!-- 操作パネル（モバイルでは最上部） -->
        <ControlPanel
          :item="auctionLiveStore.currentItem"
          :auction="auctionLiveStore.auction"
          :is-system-admin="authStore.isSystemAdmin"
          :loading="auctionLiveStore.loading"
          @start-item="handleStartItem"
          @open-price="handleOpenPrice"
          @end-item="handleEndItem"
          @end-auction="handleEndAuction"
          @cancel-auction="handleCancelAuction"
        />
        <ItemInfo :item="auctionLiveStore.currentItem" />
        <ImageSlider :media="auctionLiveStore.currentItem?.media || []" />
        <BidHistory :bids="auctionLiveStore.bids" />
        <PriceHistoryList :price-history="auctionLiveStore.priceHistory" />
        <ParticipantList :participants="auctionLiveStore.participants" />
      </div>

      <!-- デスクトップ: 3カラムレイアウト -->
      <div class="hidden lg:grid grid-cols-12 gap-6">
        <!-- 左カラム: 商品情報と画像 -->
        <div class="col-span-4 space-y-6">
          <ItemInfo :item="auctionLiveStore.currentItem" />
          <ImageSlider :media="auctionLiveStore.currentItem?.media || []" />
        </div>

        <!-- 中央カラム: 入札履歴と価格履歴 -->
        <div class="col-span-4 space-y-6">
          <BidHistory :bids="auctionLiveStore.bids" />
          <PriceHistoryList :price-history="auctionLiveStore.priceHistory" />
        </div>

        <!-- 右カラム: 操作パネルと参加者一覧 -->
        <div class="col-span-4 space-y-6">
          <ControlPanel
            :item="auctionLiveStore.currentItem"
            :auction="auctionLiveStore.auction"
            :is-system-admin="authStore.isSystemAdmin"
            :loading="auctionLiveStore.loading"
            @start-item="handleStartItem"
            @open-price="handleOpenPrice"
            @end-item="handleEndItem"
            @end-auction="handleEndAuction"
            @cancel-auction="handleCancelAuction"
          />
          <ParticipantList :participants="auctionLiveStore.participants" />
        </div>
      </div>
    </div>

    <!-- 商品一覧 (タブ形式) -->
    <Card v-if="auctionLiveStore.items.length > 0" class="p-6">
      <h3 class="text-lg font-semibold mb-4">商品一覧</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <button
          v-for="item in auctionLiveStore.items"
          :key="item.id"
          @click="handleSelectItem(item.id)"
          :class="[
            'p-4 rounded-lg border-2 transition-all text-left',
            auctionLiveStore.currentItem?.id === item.id
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300 bg-white'
          ]"
        >
          <div class="flex justify-between items-start mb-2">
            <div class="text-sm font-semibold">ロット {{ item.lot_number }}</div>
            <div
              :class="[
                'px-2 py-1 rounded-full text-xs font-medium',
                item.status === 'active' ? 'bg-green-100 text-green-700' :
                item.status === 'ended' ? 'bg-blue-100 text-blue-700' :
                'bg-gray-100 text-gray-700'
              ]"
            >
              {{ item.status === 'active' ? '進行中' : item.status === 'ended' ? '終了' : '待機中' }}
            </div>
          </div>
          <div class="text-sm font-medium truncate">{{ item.name }}</div>
          <div v-if="item.current_price" class="text-xs text-gray-500 mt-1">
            現在価格: {{ new Intl.NumberFormat('ja-JP').format(item.current_price) }} pt
          </div>
        </button>
      </div>
    </Card>
  </div>
</template>
