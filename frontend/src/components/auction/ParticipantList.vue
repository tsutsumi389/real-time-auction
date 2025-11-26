<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'

const props = defineProps({
  participants: {
    type: Array,
    default: () => [],
  },
})

const sortedParticipants = computed(() => {
  return [...props.participants].sort((a, b) => {
    if (a.status !== b.status) {
      return a.status === 'active' ? -1 : 1
    }
    return a.display_name.localeCompare(b.display_name)
  })
})

const activeCount = computed(() => {
  return props.participants.filter(p => p.status === 'active').length
})

const statusText = (status) => {
  switch (status) {
    case 'active':
      return '参加中'
    case 'suspended':
      return '一時停止'
    case 'inactive':
      return '非アクティブ'
    default:
      return status
  }
}

const statusClass = (status) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-700'
    case 'suspended':
      return 'bg-yellow-100 text-yellow-700'
    case 'inactive':
      return 'bg-gray-100 text-gray-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
}
</script>

<template>
  <Card class="p-6">
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-lg font-semibold">参加者一覧</h3>
      <div class="text-sm text-gray-500">
        アクティブ: <span class="font-semibold text-green-600">{{ activeCount }}</span> / {{ participants.length }}
      </div>
    </div>

    <div v-if="sortedParticipants.length > 0" class="space-y-2 max-h-96 overflow-y-auto">
      <div
        v-for="participant in sortedParticipants"
        :key="participant.id"
        class="flex justify-between items-center p-3 rounded-md border border-gray-200 hover:bg-gray-50 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div
            :class="[
              'w-10 h-10 rounded-full flex items-center justify-center font-semibold',
              participant.status === 'active' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
            ]"
          >
            {{ participant.display_name.charAt(0).toUpperCase() }}
          </div>
          <div>
            <div class="font-medium">{{ participant.display_name }}</div>
            <div class="text-xs text-gray-500">ID: {{ participant.id.slice(0, 8) }}</div>
          </div>
        </div>
        <div>
          <span :class="['px-2 py-1 rounded-full text-xs font-medium', statusClass(participant.status)]">
            {{ statusText(participant.status) }}
          </span>
        </div>
      </div>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
      参加者がいません
    </div>
  </Card>
</template>
