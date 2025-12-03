<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">クイックアクション</h3>

    <div class="space-y-3">
      <!-- 新規オークション作成 - 全管理者 -->
      <Button
        class="w-full justify-start"
        variant="outline"
        @click="navigateTo('/admin/auctions/create')"
      >
        <PlusIcon class="w-5 h-5 mr-2" />
        新規オークション作成
      </Button>

      <!-- 新規入札者作成 - system_adminのみ -->
      <Button
        v-if="isSystemAdmin"
        class="w-full justify-start"
        variant="outline"
        @click="navigateTo('/admin/bidders/create')"
      >
        <UserPlusIcon class="w-5 h-5 mr-2" />
        新規入札者作成
      </Button>

      <!-- ポイント付与 - system_adminのみ -->
      <Button
        v-if="isSystemAdmin"
        class="w-full justify-start"
        variant="outline"
        @click="navigateTo('/admin/points/grant')"
      >
        <CoinsIcon class="w-5 h-5 mr-2" />
        ポイント付与
      </Button>
    </div>
  </Card>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

// SVGアイコンコンポーネント
import PlusIcon from '@/components/icons/PlusIcon.vue'
import UserPlusIcon from '@/components/icons/UserPlusIcon.vue'
import CoinsIcon from '@/components/icons/CoinsIcon.vue'

const router = useRouter()
const authStore = useAuthStore()

// system_adminかどうか
const isSystemAdmin = computed(() => authStore.isSystemAdmin)

// ナビゲーション
const navigateTo = (path) => {
  router.push(path)
}
</script>
