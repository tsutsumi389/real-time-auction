/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_WS_URL: string
  // 他の環境変数もここに追加できます
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
