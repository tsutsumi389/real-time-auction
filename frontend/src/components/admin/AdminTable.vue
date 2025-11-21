<template>
  <div class="bg-white shadow rounded-lg overflow-hidden">
    <!-- ローディング表示 -->
    <div v-if="loading" class="p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">読み込み中...</p>
    </div>

    <!-- データ表示 -->
    <div v-else-if="admins.length > 0">
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
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'role')"
            >
              <div class="flex items-center gap-1">
                ロール
                <span v-if="sortField === 'role'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
            </th>
            <th
              scope="col"
              class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
              @click="$emit('sort', 'status')"
            >
              <div class="flex items-center gap-1">
                状態
                <span v-if="sortField === 'status'">
                  {{ sortOrder === 'asc' ? '↑' : '↓' }}
                </span>
              </div>
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
          <tr v-for="admin in admins" :key="admin.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ admin.id }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ admin.email }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <AdminRoleBadge :role="admin.role" />
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <AdminStatusBadge :status="admin.status" />
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(admin.created_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex gap-2">
                <button
                  @click="$emit('edit', admin.id)"
                  class="text-blue-600 hover:text-blue-900"
                  :aria-label="`${admin.email}のアカウントを編集`"
                >
                  編集
                </button>
                <button
                  @click="$emit('status-change', admin)"
                  class="text-red-600 hover:text-red-900"
                  :aria-label="`${admin.email}のアカウントを${admin.status === 'active' ? '停止' : '復活'}`"
                >
                  {{ admin.status === 'active' ? '停止' : '復活' }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- データなし表示 -->
    <div v-else class="p-8 text-center text-gray-500">
      <p>管理者が見つかりませんでした</p>
    </div>
  </div>
</template>

<script setup>
import AdminRoleBadge from './AdminRoleBadge.vue'
import AdminStatusBadge from './AdminStatusBadge.vue'

defineProps({
  admins: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  sortField: {
    type: String,
    default: 'id',
  },
  sortOrder: {
    type: String,
    default: 'asc',
  },
})

defineEmits(['sort', 'edit', 'status-change'])

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>
