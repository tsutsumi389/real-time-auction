<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">ダッシュボード</h1>
      <p class="mt-2 text-sm text-gray-600">
        システムの統計情報と最近のアクティビティ
      </p>
    </div>

    <!-- エラー表示 -->
    <Alert v-if="dashboardStore.error" variant="destructive" class="mb-6">
      <AlertTitle>エラー</AlertTitle>
      <AlertDescription>
        {{ dashboardStore.error }}
      </AlertDescription>
      <Button
        variant="outline"
        size="sm"
        class="mt-3"
        @click="handleRetry"
      >
        再試行
      </Button>
    </Alert>

    <!-- ローディング状態 -->
    <div v-if="dashboardStore.loading && !hasData" class="space-y-6">
      <!-- 統計カードのスケルトン -->
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <Skeleton v-for="i in 4" :key="i" class="h-32" />
      </div>
      <!-- コンテンツのスケルトン -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Skeleton class="h-96 lg:col-span-2" />
        <Skeleton class="h-96" />
      </div>
    </div>

    <!-- ダッシュボードコンテンツ -->
    <div v-else class="space-y-6">
      <!-- 統計カードグリッド -->
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <!-- 開催中オークション数 -->
        <StatsCard
          title="開催中オークション"
          :value="dashboardStore.stats.activeAuctions"
          :icon="ChartBarIcon"
          icon-bg-class="bg-blue-100"
          icon-color-class="text-blue-600"
        />

        <!-- 本日の入札数 -->
        <StatsCard
          title="本日の入札数"
          :value="dashboardStore.stats.todayBids"
          :icon="CursorArrowRaysIcon"
          icon-bg-class="bg-green-100"
          icon-color-class="text-green-600"
        />

        <!-- 登録入札者数 -->
        <StatsCard
          title="登録入札者数"
          :value="dashboardStore.stats.totalBidders"
          :icon="UsersIcon"
          icon-bg-class="bg-purple-100"
          icon-color-class="text-purple-600"
        />

        <!-- ポイント流通量 -->
        <StatsCard
          title="ポイント流通量"
          :value="dashboardStore.stats.totalPoints"
          :icon="CoinsIcon"
          icon-bg-class="bg-yellow-100"
          icon-color-class="text-yellow-600"
          format="points"
        />
      </div>

      <!-- コンテンツグリッド -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- メインコンテンツ（左側 2/3） -->
        <div class="space-y-6 lg:col-span-2">
          <!-- 最新の入札 -->
          <RecentBidsList :bids="dashboardStore.activities.recentBids" />

          <!-- 終了したオークション -->
          <EndedAuctionsList :auctions="dashboardStore.activities.endedAuctions" />
        </div>

        <!-- サイドバー（右側 1/3） -->
        <div class="space-y-6">
          <!-- クイックアクション -->
          <QuickActions />

          <!-- 新規入札者（system_adminのみ） -->
          <NewBiddersList
            v-if="authStore.isSystemAdmin"
            :bidders="dashboardStore.activities.newBidders"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useDashboardStore } from '@/stores/dashboardStore'
import { useAuthStore } from '@/stores/auth'

// UIコンポーネント
import Alert from '@/components/ui/Alert.vue'
import AlertTitle from '@/components/ui/AlertTitle.vue'
import AlertDescription from '@/components/ui/AlertDescription.vue'
import Button from '@/components/ui/Button.vue'
import Skeleton from '@/components/ui/Skeleton.vue'

// ダッシュボードコンポーネント
import StatsCard from '@/components/admin/StatsCard.vue'
import RecentBidsList from '@/components/admin/RecentBidsList.vue'
import NewBiddersList from '@/components/admin/NewBiddersList.vue'
import EndedAuctionsList from '@/components/admin/EndedAuctionsList.vue'
import QuickActions from '@/components/admin/QuickActions.vue'

// アイコン
import ChartBarIcon from '@/components/icons/ChartBarIcon.vue'
import CursorArrowRaysIcon from '@/components/icons/CursorArrowRaysIcon.vue'
import UsersIcon from '@/components/icons/UsersIcon.vue'
import CoinsIcon from '@/components/icons/CoinsIcon.vue'

const dashboardStore = useDashboardStore()
const authStore = useAuthStore()

// データが存在するかチェック
const hasData = computed(() => {
  return (
    dashboardStore.stats.activeAuctions > 0 ||
    dashboardStore.stats.todayBids > 0 ||
    dashboardStore.stats.totalBidders > 0 ||
    dashboardStore.stats.totalPoints > 0 ||
    dashboardStore.activities.recentBids.length > 0 ||
    dashboardStore.activities.endedAuctions.length > 0 ||
    dashboardStore.activities.newBidders.length > 0
  )
})

// 再試行ハンドラー
const handleRetry = async () => {
  dashboardStore.clearError()
  await dashboardStore.fetchAll()
}

// マウント時にデータ取得
onMounted(async () => {
  await dashboardStore.fetchAll()
})
</script>
