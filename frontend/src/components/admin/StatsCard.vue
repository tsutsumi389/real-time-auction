<template>
  <Card
    class="p-6 hover:shadow-lg transition-shadow duration-200 cursor-default"
  >
    <!-- アイコンとタイトル -->
    <div class="flex items-center justify-between mb-4">
      <div class="flex items-center gap-3">
        <div
          :class="[
            'p-3 rounded-lg',
            iconBgClass
          ]"
        >
          <component
            :is="icon"
            :class="['w-6 h-6', iconColorClass]"
          />
        </div>
        <h3 class="text-sm font-medium text-gray-600">{{ title }}</h3>
      </div>
    </div>

    <!-- 値の表示 -->
    <div class="mt-2">
      <p class="text-3xl font-bold text-gray-900">
        {{ formattedValue }}
      </p>
      <p v-if="subtitle" class="text-xs text-gray-500 mt-1">
        {{ subtitle }}
      </p>
    </div>
  </Card>
</template>

<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'

const props = defineProps({
  /** カードのタイトル */
  title: {
    type: String,
    required: true,
  },
  /** 表示する値 */
  value: {
    type: Number,
    required: true,
  },
  /** アイコンコンポーネント（SVGコンポーネント） */
  icon: {
    type: Object,
    default: null,
  },
  /** アイコンの背景色クラス */
  iconBgClass: {
    type: String,
    default: 'bg-blue-100',
  },
  /** アイコンの色クラス */
  iconColorClass: {
    type: String,
    default: 'text-blue-600',
  },
  /** サブタイトル（オプション） */
  subtitle: {
    type: String,
    default: '',
  },
  /** 数値フォーマット（"number" | "points"） */
  format: {
    type: String,
    default: 'number',
    validator: (value) => ['number', 'points'].includes(value),
  },
})

// 値のフォーマット
const formattedValue = computed(() => {
  if (props.format === 'points') {
    return `${props.value.toLocaleString()} pt`
  }
  return props.value.toLocaleString()
})
</script>
