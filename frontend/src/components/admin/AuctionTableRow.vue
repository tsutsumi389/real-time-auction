<template>
  <tr class="hover:bg-gray-50">
    <!-- ID -->
    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-500">
      {{ auction.id.substring(0, 8) }}...
    </td>

    <!-- タイトル -->
    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
      <div class="flex items-center gap-2">
        {{ auction.title }}
        <!-- 商品数バッジ -->
        <span
          v-if="auction.item_count === 0"
          class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600"
        >
          🏷️ 0件
        </span>
        <span
          v-else
          class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700"
        >
          🏷️ {{ auction.item_count }}件
        </span>
      </div>
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
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <button
            class="inline-flex items-center justify-center w-8 h-8 rounded hover:bg-gray-100 transition-colors"
            title="アクション"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-5 w-5 text-gray-500"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
            </svg>
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-48">
          <!-- 編集（pendingのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'pending'"
            @select="$emit('edit', auction.id)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
            </svg>
            編集
          </DropdownMenuItem>

          <!-- 商品追加（pendingのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'pending'"
            @select="handleAddItem(auction.id)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
            商品追加
          </DropdownMenuItem>

          <!-- セパレータ -->
          <DropdownMenuSeparator v-if="auction.status === 'pending'" />

          <!-- ライブ表示（activeのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'active'"
            @select="$emit('view-live', auction.id)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
              <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
            </svg>
            ライブ表示
          </DropdownMenuItem>

          <!-- 詳細（常に表示） -->
          <DropdownMenuItem @select="$emit('view-details', auction.id)">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
              <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
            </svg>
            詳細
          </DropdownMenuItem>

          <!-- セパレータ -->
          <DropdownMenuSeparator v-if="auction.status === 'pending' || auction.status === 'active'" />

          <!-- 公開（pendingのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'pending'"
            @select="$emit('start', auction)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
            </svg>
            公開
          </DropdownMenuItem>

          <!-- 終了（activeのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'active'"
            @select="$emit('end', auction)"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8 7a1 1 0 00-1 1v4a1 1 0 001 1h4a1 1 0 001-1V8a1 1 0 00-1-1H8z" clip-rule="evenodd" />
            </svg>
            終了
          </DropdownMenuItem>

          <!-- 中止（activeのみ、system_adminのみ） -->
          <DropdownMenuItem
            v-if="auction.status === 'active' && isSystemAdmin"
            @select="$emit('cancel', auction)"
            class="text-red-600 focus:text-red-600"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4 mr-2"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            中止
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </td>
  </tr>
</template>

<script setup>
import { useRouter } from 'vue-router'
import AuctionStatusBadge from './AuctionStatusBadge.vue'
import DropdownMenu from '../ui/DropdownMenu.vue'
import DropdownMenuTrigger from '../ui/DropdownMenuTrigger.vue'
import DropdownMenuContent from '../ui/DropdownMenuContent.vue'
import DropdownMenuItem from '../ui/DropdownMenuItem.vue'
import DropdownMenuSeparator from '../ui/DropdownMenuSeparator.vue'

const router = useRouter()

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

function handleAddItem(auctionId) {
  router.push(`/admin/auctions/${auctionId}/edit?tab=create`)
}
</script>
