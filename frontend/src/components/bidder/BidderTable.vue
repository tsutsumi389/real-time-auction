<template>
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <!-- ローディング表示 -->
    <div v-if="loading" class="p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">読み込み中...</p>
    </div>

    <!-- データ表示 -->
    <div v-else-if="bidders.length > 0">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'id')"
            >
              <div class="flex items-center gap-1">
                ID
                <span v-if="sortField === 'id'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'email')"
            >
              <div class="flex items-center gap-1">
                メールアドレス
                <span v-if="sortField === 'email'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              表示名
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'points')"
            >
              <div class="flex items-center gap-1">
                ポイント
                <span v-if="sortField === 'points'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
            >
              状態
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'created_at')"
            >
              <div class="flex items-center gap-1">
                作成日
                <span v-if="sortField === 'created_at'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              アクション
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="bidder in bidders" :key="bidder.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              <span :title="bidder.id" class="cursor-help">
                {{ formatId(bidder.id) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ bidder.email }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ bidder.display_name || '（未設定）' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
              {{ formatPoints(bidder.points) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <AdminStatusBadge :status="bidder.status" />
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(bidder.created_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex gap-2">
                <button
                  @click="$emit('edit', bidder.id)"
                  class="text-gray-600 hover:text-gray-900"
                  :aria-label="`${bidder.email}のアカウントを編集`"
                >
                  詳細
                </button>
                <button
                  v-if="bidder.status !== 'deleted'"
                  @click="$emit('grant-points', bidder)"
                  class="text-green-600 hover:text-green-900"
                  :aria-label="`${bidder.email}にポイントを付与`"
                >
                  pt付与
                </button>
                <button
                  @click="$emit('show-history', bidder)"
                  class="text-purple-600 hover:text-purple-900"
                  :aria-label="`${bidder.email}のポイント履歴を表示`"
                >
                  履歴
                </button>
                <button
                  v-if="bidder.status !== 'deleted'"
                  @click="$emit('status-change', bidder)"
                  :class="bidder.status === 'active' ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'"
                  :aria-label="`${bidder.email}のアカウントを${bidder.status === 'active' ? '停止' : '復活'}`"
                >
                  {{ bidder.status === 'active' ? '停止' : '復活' }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- データなし表示 -->
    <div v-else class="p-8 text-center text-gray-500">
      <p>入札者が見つかりませんでした</p>
    </div>
  </div>
</template>

<script setup>
import AdminStatusBadge from '../admin/AdminStatusBadge.vue'

defineProps({
  bidders: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  sortField: {
    type: String,
    default: 'created_at',
  },
  sortOrder: {
    type: String,
    default: 'asc',
  },
})

defineEmits(['sort', 'edit', 'grant-points', 'show-history', 'status-change'])

// UUID先頭8文字のみ表示
function formatId(uuid) {
  return uuid.substring(0, 8)
}

// ポイントを3桁区切りで表示
function formatPoints(points) {
  return (points || 0).toLocaleString('ja-JP')
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>
