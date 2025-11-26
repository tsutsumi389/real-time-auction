<script setup>
import { ref, computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Dialog from '@/components/ui/Dialog.vue'

const props = defineProps({
  item: {
    type: Object,
    default: null,
  },
  auction: {
    type: Object,
    default: null,
  },
  isSystemAdmin: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits([
  'start-item',
  'open-price',
  'end-item',
  'end-auction',
  'cancel-auction',
])

const priceInput = ref('')
const showConfirm = ref(null)

const canStartItem = computed(() => {
  return props.item && props.item.status === 'pending'
})

const canOpenPrice = computed(() => {
  return props.item && props.item.status === 'active'
})

const canEndItem = computed(() => {
  return props.item && props.item.status === 'active'
})

const canEndAuction = computed(() => {
  return props.auction && props.auction.status === 'active'
})

const canCancelAuction = computed(() => {
  return props.isSystemAdmin && props.auction && ['pending', 'active'].includes(props.auction.status)
})

const isPriceValid = computed(() => {
  const price = parseInt(priceInput.value)
  return !isNaN(price) && price > 0
})

function handleStartItem() {
  showConfirm.value = 'start-item'
}

function handleOpenPrice() {
  if (!isPriceValid.value) return
  showConfirm.value = 'open-price'
}

function handleEndItem() {
  showConfirm.value = 'end-item'
}

function handleEndAuction() {
  showConfirm.value = 'end-auction'
}

function handleCancelAuction() {
  showConfirm.value = 'cancel-auction'
}

function confirmAction() {
  const action = showConfirm.value

  switch (action) {
    case 'start-item':
      emit('start-item', props.item.id)
      break
    case 'open-price':
      emit('open-price', props.item.id, parseInt(priceInput.value))
      priceInput.value = ''
      break
    case 'end-item':
      emit('end-item', props.item.id)
      break
    case 'end-auction':
      emit('end-auction', props.auction.id)
      break
    case 'cancel-auction':
      emit('cancel-auction', props.auction.id)
      break
  }

  showConfirm.value = null
}

function cancelAction() {
  showConfirm.value = null
}

const confirmTitle = computed(() => {
  switch (showConfirm.value) {
    case 'start-item':
      return '商品開始の確認'
    case 'open-price':
      return '価格開示の確認'
    case 'end-item':
      return '商品終了の確認'
    case 'end-auction':
      return 'オークション終了の確認'
    case 'cancel-auction':
      return '緊急停止の確認'
    default:
      return '確認'
  }
})

const confirmMessage = computed(() => {
  switch (showConfirm.value) {
    case 'start-item':
      return '商品を開始しますか?'
    case 'open-price':
      return `価格 ${priceInput.value} pt を開示しますか?`
    case 'end-item':
      return '商品を終了しますか? 最高入札者が落札者となります。'
    case 'end-auction':
      return 'オークションを終了しますか?'
    case 'cancel-auction':
      return 'オークションを緊急停止しますか? すべての入札が取り消され、ポイントが返金されます。'
    default:
      return ''
  }
})
</script>

<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold mb-4">操作パネル</h3>

    <!-- 確認ダイアログ（モーダル） -->
    <Dialog :open="showConfirm !== null" @update:open="val => !val && cancelAction()" :title="confirmTitle">
      <p class="text-sm text-gray-600">{{ confirmMessage }}</p>
      <template #footer="{ close }">
        <Button @click="confirmAction" variant="default" size="sm" :disabled="loading">
          確認
        </Button>
        <Button @click="close" variant="outline" size="sm" :disabled="loading">
          キャンセル
        </Button>
      </template>
    </Dialog>

    <!-- 商品操作 -->
    <div class="space-y-4">
      <div>
        <Label for="price-input" class="mb-2">価格開示</Label>
        <div class="flex gap-2">
          <Input
            id="price-input"
            v-model="priceInput"
            type="number"
            placeholder="価格を入力"
            :disabled="!canOpenPrice || loading"
            @keyup.enter="handleOpenPrice"
          />
          <Button
            @click="handleOpenPrice"
            :disabled="!canOpenPrice || !isPriceValid || loading"
          >
            開示
          </Button>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Button
          @click="handleStartItem"
          :disabled="!canStartItem || loading"
          variant="default"
          class="w-full"
        >
          商品開始
        </Button>

        <Button
          @click="handleEndItem"
          :disabled="!canEndItem || loading"
          variant="default"
          class="w-full"
        >
          商品終了
        </Button>
      </div>
    </div>

    <!-- オークション操作 -->
    <div class="mt-6 pt-6 border-t space-y-2">
      <Button
        @click="handleEndAuction"
        :disabled="!canEndAuction || loading"
        variant="outline"
        class="w-full"
      >
        オークション終了
      </Button>

      <Button
        v-if="isSystemAdmin"
        @click="handleCancelAuction"
        :disabled="!canCancelAuction || loading"
        variant="destructive"
        class="w-full"
      >
        緊急停止
      </Button>
    </div>
  </Card>
</template>
