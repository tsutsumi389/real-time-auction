<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <h2 class="text-lg font-semibold text-gray-900 mb-4">ポイント残高</h2>

    <!-- Points Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
      <!-- Total Points -->
      <div class="text-center p-4 bg-blue-50 rounded-lg border border-blue-100">
        <div class="text-sm text-gray-600 mb-1">合計ポイント</div>
        <div class="text-3xl font-bold text-blue-600">
          {{ formatNumber(points.total) }}
        </div>
        <div class="text-xs text-gray-500 mt-1">全ポイント</div>
      </div>

      <!-- Available Points -->
      <div class="text-center p-4 bg-green-50 rounded-lg border border-green-100">
        <div class="text-sm text-gray-600 mb-1">利用可能</div>
        <div class="text-3xl font-bold text-green-600">
          {{ formatNumber(points.available) }}
        </div>
        <div class="text-xs text-gray-500 mt-1">入札可能</div>
      </div>

      <!-- Reserved Points -->
      <div class="text-center p-4 bg-yellow-50 rounded-lg border border-yellow-100">
        <div class="text-sm text-gray-600 mb-1">予約済み</div>
        <div class="text-3xl font-bold text-yellow-600">
          {{ formatNumber(points.reserved) }}
        </div>
        <div class="text-xs text-gray-500 mt-1">入札中</div>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="space-y-2">
      <div class="flex justify-between text-sm text-gray-600">
        <span>ポイント利用状況</span>
        <span>{{ usagePercentage }}%</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
        <div class="h-full flex">
          <!-- Available portion (green) -->
          <div
            :style="{ width: availablePercentage + '%' }"
            class="bg-green-500 transition-all duration-500"
          ></div>
          <!-- Reserved portion (yellow) -->
          <div
            :style="{ width: reservedPercentage + '%' }"
            class="bg-yellow-500 transition-all duration-500"
          ></div>
        </div>
      </div>
      <div class="flex justify-between text-xs text-gray-500">
        <span>
          <span class="inline-block w-2 h-2 bg-green-500 rounded-full mr-1"></span>
          利用可能: {{ formatNumber(points.available) }}
        </span>
        <span>
          <span class="inline-block w-2 h-2 bg-yellow-500 rounded-full mr-1"></span>
          予約済み: {{ formatNumber(points.reserved) }}
        </span>
      </div>
    </div>

    <!-- Warning for low points -->
    <div
      v-if="points.available < 10000 && points.available > 0"
      class="mt-4 p-3 bg-orange-50 border border-orange-200 rounded-lg"
    >
      <p class="text-sm text-orange-800">
        ⚠️ 利用可能ポイントが少なくなっています
      </p>
    </div>

    <!-- No points warning -->
    <div
      v-if="points.available === 0 && points.reserved > 0"
      class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg"
    >
      <p class="text-sm text-yellow-800">
        💡 すべてのポイントが予約済みです。新たに入札するには、他の商品で落札されないか、商品が終了するまでお待ちください。
      </p>
    </div>

    <!-- Zero points warning -->
    <div
      v-if="points.available === 0 && points.reserved === 0 && points.total === 0"
      class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg"
    >
      <p class="text-sm text-red-800">
        ❌ ポイントがありません。管理者にポイント付与を依頼してください。
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  points: {
    type: Object,
    required: true,
    validator: (value) => {
      return (
        typeof value.total === 'number' &&
        typeof value.available === 'number' &&
        typeof value.reserved === 'number'
      )
    },
  },
})

// Computed properties for progress bar
const availablePercentage = computed(() => {
  if (props.points.total === 0) return 0
  return Math.round((props.points.available / props.points.total) * 100)
})

const reservedPercentage = computed(() => {
  if (props.points.total === 0) return 0
  return Math.round((props.points.reserved / props.points.total) * 100)
})

const usagePercentage = computed(() => {
  return availablePercentage.value + reservedPercentage.value
})

// Format number with comma separator
function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}
</script>
