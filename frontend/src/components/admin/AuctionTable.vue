<template>
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <!-- ローディング表示 -->
    <div v-if="loading" class="p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">読み込み中...</p>
    </div>

    <!-- データ表示 -->
    <div v-else-if="auctions.length > 0">
      <div class="overflow-x-auto">
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
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                タイトル
              </th>
              <th
                scope="col"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                説明
              </th>
              <th
                scope="col"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                状態
              </th>
              <th
                scope="col"
                class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                商品数
              </th>
              <th
                scope="col"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                @click="$emit('sort', 'created_at')"
              >
                <div class="flex items-center gap-1">
                  作成日時
                  <span v-if="sortField === 'created_at'">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </div>
              </th>
              <th
                scope="col"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                @click="$emit('sort', 'updated_at')"
              >
                <div class="flex items-center gap-1">
                  更新日時
                  <span v-if="sortField === 'updated_at'">
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
            <AuctionTableRow
              v-for="auction in auctions"
              :key="auction.id"
              :auction="auction"
              :is-system-admin="isSystemAdmin"
              @view-details="$emit('view-details', $event)"
              @view-live="$emit('view-live', $event)"
              @start="$emit('start', $event)"
              @end="$emit('end', $event)"
              @cancel="$emit('cancel', $event)"
              @edit="$emit('edit', $event)"
            />
          </tbody>
        </table>
      </div>
    </div>

    <!-- データなし表示 -->
    <div v-else class="p-8 text-center text-gray-500">
      <p>該当するオークションが見つかりませんでした</p>
    </div>
  </div>
</template>

<script setup>
import AuctionTableRow from './AuctionTableRow.vue'

defineProps({
  auctions: {
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
    default: 'desc',
  },
  isSystemAdmin: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['sort', 'view-details', 'view-live', 'start', 'end', 'cancel', 'edit'])
</script>
