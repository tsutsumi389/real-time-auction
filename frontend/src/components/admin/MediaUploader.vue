<script setup>
/**
 * MediaUploader Component
 * ドラッグ&ドロップ対応のメディアアップロードコンポーネント
 */
import { ref, computed } from 'vue'
import {
  uploadItemMedia,
  validateFile,
  formatFileSize,
  getMediaTypeFromFile,
} from '@/services/mediaApi'

const props = defineProps({
  itemId: {
    type: String,
    required: true,
  },
  maxImages: {
    type: Number,
    default: 10,
  },
  maxVideos: {
    type: Number,
    default: 3,
  },
  currentImageCount: {
    type: Number,
    default: 0,
  },
  currentVideoCount: {
    type: Number,
    default: 0,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['upload-success', 'upload-error'])

// ローカル状態
const isDragOver = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadQueue = ref([])
const currentUploadIndex = ref(0)
const error = ref('')

// 残り枠数の計算
const remainingImages = computed(() => props.maxImages - props.currentImageCount)
const remainingVideos = computed(() => props.maxVideos - props.currentVideoCount)
const canUploadImage = computed(() => remainingImages.value > 0 && !props.disabled)
const canUploadVideo = computed(() => remainingVideos.value > 0 && !props.disabled)

// ドラッグ&ドロップイベントハンドラー
function handleDragEnter(e) {
  e.preventDefault()
  if (!props.disabled) {
    isDragOver.value = true
  }
}

function handleDragOver(e) {
  e.preventDefault()
  if (!props.disabled) {
    isDragOver.value = true
  }
}

function handleDragLeave(e) {
  e.preventDefault()
  isDragOver.value = false
}

function handleDrop(e) {
  e.preventDefault()
  isDragOver.value = false

  if (props.disabled) return

  const files = Array.from(e.dataTransfer.files)
  handleFiles(files)
}

// ファイル入力のchangeイベント
function handleFileInput(e) {
  const files = Array.from(e.target.files)
  handleFiles(files)
  // 入力をリセットして同じファイルを再選択可能に
  e.target.value = ''
}

// ファイル処理
async function handleFiles(files) {
  error.value = ''

  // ファイルのバリデーション
  const validFiles = []
  let imageCount = 0
  let videoCount = 0

  for (const file of files) {
    const validationError = validateFile(file)
    if (validationError) {
      error.value = validationError
      continue
    }

    const mediaType = getMediaTypeFromFile(file)

    if (mediaType === 'image') {
      if (imageCount + props.currentImageCount >= props.maxImages) {
        error.value = `画像は最大${props.maxImages}枚までです`
        continue
      }
      imageCount++
    } else if (mediaType === 'video') {
      if (videoCount + props.currentVideoCount >= props.maxVideos) {
        error.value = `動画は最大${props.maxVideos}本までです`
        continue
      }
      videoCount++
    }

    validFiles.push({
      file,
      mediaType,
      status: 'pending',
      progress: 0,
    })
  }

  if (validFiles.length === 0) {
    return
  }

  // アップロードキューに追加
  uploadQueue.value = validFiles
  currentUploadIndex.value = 0
  uploading.value = true

  // 順番にアップロード
  for (let i = 0; i < validFiles.length; i++) {
    currentUploadIndex.value = i
    uploadQueue.value[i].status = 'uploading'

    try {
      const result = await uploadItemMedia(
        props.itemId,
        validFiles[i].file,
        validFiles[i].mediaType,
        (progress) => {
          uploadProgress.value = progress
          uploadQueue.value[i].progress = progress
        }
      )

      uploadQueue.value[i].status = 'success'
      emit('upload-success', result)
    } catch (err) {
      uploadQueue.value[i].status = 'error'
      error.value = err.message || 'アップロードに失敗しました'
      emit('upload-error', err)
    }
  }

  // アップロード完了
  uploading.value = false
  uploadProgress.value = 0

  // 少し待ってからキューをクリア
  setTimeout(() => {
    uploadQueue.value = []
  }, 2000)
}

// ファイル選択ダイアログを開く
function openFileDialog() {
  if (!props.disabled) {
    document.getElementById('file-input').click()
  }
}
</script>

<template>
  <div class="media-uploader">
    <!-- 残り枠数の表示 -->
    <div class="mb-3 flex gap-4 text-sm text-gray-600">
      <span>
        画像: <span :class="{ 'text-red-500': remainingImages <= 0 }">{{ currentImageCount }}/{{ maxImages }}</span>
      </span>
      <span>
        動画: <span :class="{ 'text-red-500': remainingVideos <= 0 }">{{ currentVideoCount }}/{{ maxVideos }}</span>
      </span>
    </div>

    <!-- ドラッグ&ドロップエリア -->
    <div
      :class="[
        'relative border-2 border-dashed rounded-lg p-6 text-center transition-all duration-200',
        isDragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400',
        disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer',
      ]"
      @click="openFileDialog"
      @dragenter="handleDragEnter"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
    >
      <!-- アップロードアイコン -->
      <div class="flex flex-col items-center">
        <svg
          v-if="!uploading"
          class="w-12 h-12 text-gray-400 mb-3"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
          />
        </svg>

        <!-- アップロード中のスピナー -->
        <div v-else class="mb-3">
          <svg
            class="w-12 h-12 text-blue-500 animate-spin"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        </div>

        <!-- テキスト -->
        <p v-if="!uploading" class="text-gray-600">
          <span class="font-medium text-blue-600">クリックしてファイルを選択</span>
          <br />
          またはドラッグ&ドロップ
        </p>
        <p v-else class="text-gray-600">
          アップロード中...
          <span class="font-medium">{{ uploadProgress }}%</span>
        </p>

        <!-- 対応形式 -->
        <p class="mt-2 text-xs text-gray-500">
          画像: JPEG, PNG, WebP, GIF（最大5MB）<br />
          動画: MP4, MOV, AVI（最大100MB）
        </p>
      </div>

      <!-- プログレスバー -->
      <div
        v-if="uploading"
        class="absolute bottom-0 left-0 h-1 bg-blue-500 transition-all duration-300"
        :style="{ width: `${uploadProgress}%` }"
      />

      <!-- 隠れたファイル入力 -->
      <input
        id="file-input"
        type="file"
        multiple
        accept="image/jpeg,image/png,image/webp,image/gif,video/mp4,video/quicktime,video/x-msvideo"
        class="hidden"
        :disabled="disabled"
        @change="handleFileInput"
      />
    </div>

    <!-- アップロードキュー -->
    <div v-if="uploadQueue.length > 0" class="mt-4 space-y-2">
      <div
        v-for="(item, index) in uploadQueue"
        :key="index"
        class="flex items-center gap-3 p-2 bg-gray-50 rounded-lg"
      >
        <!-- ステータスアイコン -->
        <div class="flex-shrink-0">
          <!-- 待機中 -->
          <svg
            v-if="item.status === 'pending'"
            class="w-5 h-5 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <!-- アップロード中 -->
          <svg
            v-else-if="item.status === 'uploading'"
            class="w-5 h-5 text-blue-500 animate-spin"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          <!-- 成功 -->
          <svg
            v-else-if="item.status === 'success'"
            class="w-5 h-5 text-green-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
            />
          </svg>
          <!-- エラー -->
          <svg
            v-else-if="item.status === 'error'"
            class="w-5 h-5 text-red-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </div>

        <!-- ファイル情報 -->
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-gray-900 truncate">
            {{ item.file.name }}
          </p>
          <p class="text-xs text-gray-500">
            {{ formatFileSize(item.file.size) }}
            <span v-if="item.status === 'uploading'" class="ml-2">{{ item.progress }}%</span>
          </p>
        </div>

        <!-- プログレスバー -->
        <div v-if="item.status === 'uploading'" class="w-24 h-1.5 bg-gray-200 rounded-full overflow-hidden">
          <div
            class="h-full bg-blue-500 transition-all duration-300"
            :style="{ width: `${item.progress}%` }"
          />
        </div>
      </div>
    </div>

    <!-- エラーメッセージ -->
    <div v-if="error" class="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg">
      <div class="flex items-start gap-2">
        <svg
          class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <p class="text-sm text-red-700">{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.media-uploader {
  width: 100%;
}
</style>
