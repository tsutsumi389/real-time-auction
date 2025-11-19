# Frontend

Vue.js 3 + Vite によるフロントエンドアプリケーション

## 構造

```
frontend/
├── src/
│   ├── assets/           # 静的アセット（CSS、画像など）
│   ├── components/       # 再利用可能なコンポーネント
│   ├── views/            # ページコンポーネント
│   ├── router/           # Vue Router設定
│   ├── stores/           # Pinia ストア（状態管理）
│   ├── services/         # API通信サービス
│   ├── utils/            # ユーティリティ関数
│   ├── App.vue           # ルートコンポーネント
│   └── main.js           # エントリーポイント
├── public/               # 公開静的ファイル
├── index.html            # HTMLテンプレート
├── vite.config.js        # Vite設定
├── package.json          # npm依存関係
├── Dockerfile            # 本番用Dockerfile
├── Dockerfile.dev        # 開発用Dockerfile
└── nginx.conf            # Nginx設定（本番用）
```

## 開発

### ローカル開発（Dockerなし）

```bash
# 依存関係のインストール
npm install

# 開発サーバー起動
npm run dev

# ビルド
npm run build
```

### Docker開発環境

プロジェクトルートから:

```bash
# 起動
make up

# ログ確認
make logs-frontend

# アクセス
open http://localhost
```

## 実装予定

- [ ] 認証画面（ログイン・登録）
- [ ] システム管理画面
- [ ] 主催者管理画面
- [ ] オークション一覧・詳細画面
- [ ] リアルタイム入札UI
- [ ] WebSocket接続管理
- [ ] ユーザープロフィール
- [ ] 入札履歴
