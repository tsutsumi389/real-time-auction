<template>
  <tr class="hover:bg-gray-50">
    <!-- ID -->
    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-500">
      {{ auction.id.substring(0, 8) }}...
    </td>

    <!-- タイトル -->
    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
      {{ auction.title }}
    </td>

    <!-- 説明 -->
    <td class="px-6 py-4 text-sm text-gray-500">
      {{ truncateDescription(auction.description) }}
    </td>

    <!-- 状態 -->
    <td class="px-6 py-4 whitespace-nowrap">
      <AuctionStatusBadge :status="auction.status" />
    </td>

    <!-- 商品数 -->
    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-center">
      {{ auction.item_count }}
    </td>

    <!-- 作成日時 -->
    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
      {{ formatDate(auction.created_at) }}
    </td>

    <!-- 更新日時 -->
    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
      {{ formatDate(auction.updated_at) }}
    </td>

    <!-- アクション -->
    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
      <div class="flex gap-2">
        <!-- 詳細ボタン（常に表示） -->
        <button
          @click="$emit('view-details', auction.id)"
          class="text-blue-600 hover:text-blue-900"
          title="詳細"
        >
          詳細
        </button>

        <!-- 編集ボタン（pendingのみ） -->
        <button
          v-if="auction.status === 'pending'"
          @click="$emit('edit', auction.id)"
          class="text-indigo-600 hover:text-indigo-900"
          title="編集"
        >
          編集
        </button>

        <!-- 開催中ボタン（activeのみ） -->
        <button
          v-if="auction.status === 'active'"
          @click="$emit('view-live', auction.id)"
          class="text-green-600 hover:text-green-900"
          title="開催中"
        >
          開催中
        </button>

        <!-- 開始ボタン（pendingのみ） -->
        <button
          v-if="auction.status === 'pending'"
          @click="$emit('start', auction)"
          class="text-green-600 hover:text-green-900"
          title="公開"
        >
          公開
        </button>

        <!-- 終了ボタン（activeのみ） -->
        <button
          v-if="auction.status === 'active'"
          @click="$emit('end', auction)"
          class="text-blue-600 hover:text-blue-900"
          title="終了"
        >
          終了
        </button>

        <!-- 中止ボタン（activeのみ、system_adminのみ） -->
        <button
          v-if="auction.status === 'active' && isSystemAdmin"
          @click="$emit('cancel', auction)"
          class="text-red-600 hover:text-red-900"
          title="中止"
        >
          中止
        </button>
      </div>
    </td>
  </tr>
</template>

<script setup>
import AuctionStatusBadge from './AuctionStatusBadge.vue'

const props = defineProps({
  auction: {
    type: Object,
    required: true,
  },
  isSystemAdmin: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['view-details', 'view-live', 'start', 'end', 'cancel', 'edit'])

function truncateDescription(description) {
  if (!description) return '-'
  if (description.length <= 50) return description
  return description.substring(0, 50) + '...'
}

function formatDate(dateString) {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>
