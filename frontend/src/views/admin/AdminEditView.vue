<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">管理者編集</h1>
      <p class="mt-2 text-sm text-gray-600">管理者アカウントの編集（未実装）</p>
    </div>

    <!-- プレースホルダーメッセージ -->
    <div class="bg-white shadow rounded-lg p-6">
      <div class="text-center py-12">
        <p class="text-gray-500 mb-4">管理者編集画面は実装中です</p>
        <button
          @click="$router.push({ name: 'admin-list' })"
          class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
        >
          管理者一覧に戻る
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// プレースホルダー実装
// TODO: 管理者編集機能を実装する
</script>
