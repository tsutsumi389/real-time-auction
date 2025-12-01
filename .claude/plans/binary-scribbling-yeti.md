# MinIO画像保存機能 実装計画

## 背景と決定事項

### 現在の状況
- **データベーススキーマ**: `item_media` テーブルは実装済み（URL、サムネイルURL、表示順序を保存）
- **バックエンド実装**: 画像アップロード機能は完全未実装
- **Docker環境**: 基本的なボリュームマウント設定済み
- **本番環境**: AWS S3 または GCP Cloud Storage を使用予定

### 決定事項
- **ストレージ**: MinIO（S3互換）をローカル環境で使用
- **画像処理**: リサイズ（最大1920x1080）、JPEG変換、サムネイル生成（300x300）を実装
  - **画像フォーマット変更（2025-12-01）**: 当初WebPを予定していたが、CGO依存の問題によりJPEG形式に変更
  - WebPライブラリ（chai2010/webp、kolesa-team/go-webp）はlibwebp（C言語実装）をCGO経由で呼び出すため、Docker環境でビルドエラーが発生
  - `github.com/disintegration/imaging` を使用したJPEG形式（品質80）で高品質を維持
  - 純粋なGo実装により、クロスコンパイルとDocker環境での安定性を確保
- **サムネイル**: アップロード時に自動生成
- **ファイルサイズ制限**: 画像5MB、動画100MB
- **MinIO**: AGPLv3ライセンスで完全無料、S3互換APIを使用

### MinIOを選択した理由
1. **本番とコードが完全に同じ**: S3 SDKをそのまま使えるため環境差分によるバグを防止
2. **開発体験**: Web UIなしでもCLIとAPIで十分な管理が可能
3. **無料**: AGPLv3ライセンス、内部開発に最適

---

## 1. アーキテクチャ設計

### レイヤー構造
```
[Client] → [Nginx] → [API Handler]
                          ↓
                   [MediaHandler]
                          ↓
              ┌───────────┴───────────┐
              ↓                       ↓
      [ImageProcessor]        [StorageService]
        (リサイズ/JPEG)         (MinIO/S3/GCS)
              ↓                       ↓
      [ItemMediaRepository]           |
              ↓                       |
        [PostgreSQL] ←────────────────┘
```

### ファイルアップロードフロー
1. JWT認証・ファイルサイズ検証（middleware）
2. MIMEタイプ検証・ファイル名サニタイズ
3. 画像処理（リサイズ → JPEG変換 → サムネイル生成）
4. ストレージアップロード（オリジナル + サムネイル）
5. DB保存（トランザクション）
6. 一時ファイル削除
7. エラー時は完全ロールバック

---

## 2. MinIO環境構築

### docker-compose.yml 追加内容

```yaml
services:
  minio:
    image: minio/minio:latest
    container_name: auction-minio
    ports:
      - "9000:9000"  # API
      - "9001:9001"  # Console
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY:-minioadmin}
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    networks:
      - auction-network
    healthcheck:
      test: ["CMD", "mc", "ready", "local"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio-init:
    image: minio/mc:latest
    container_name: auction-minio-init
    depends_on:
      minio:
        condition: service_healthy
    environment:
      MINIO_ENDPOINT: minio:9000
      MINIO_ACCESS_KEY: ${MINIO_ACCESS_KEY:-minioadmin}
      MINIO_SECRET_KEY: ${MINIO_SECRET_KEY:-minioadmin}
      MINIO_BUCKET: ${MINIO_BUCKET:-auction-media}
    networks:
      - auction-network
    entrypoint: >
      /bin/sh -c "
      mc alias set local http://$$MINIO_ENDPOINT $$MINIO_ACCESS_KEY $$MINIO_SECRET_KEY;
      mc mb local/$$MINIO_BUCKET --ignore-existing;
      mc anonymous set download local/$$MINIO_BUCKET;
      exit 0;
      "

volumes:
  minio_data:
    driver: local
```

### .env 追加項目

```bash
# Storage Configuration
STORAGE_TYPE=minio

# MinIO Settings (Local)
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=auction-media
MINIO_USE_SSL=false
MINIO_PUBLIC_URL=http://localhost:9000

# AWS S3 Settings (Production)
S3_REGION=ap-northeast-1
S3_BUCKET=auction-media-prod
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=

# Upload Limits
MAX_IMAGE_SIZE=5242880       # 5MB
MAX_VIDEO_SIZE=104857600     # 100MB
IMAGE_MAX_WIDTH=1920
IMAGE_MAX_HEIGHT=1080
THUMBNAIL_SIZE=300

# Allowed MIME Types
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp,image/gif
ALLOWED_VIDEO_TYPES=video/mp4,video/quicktime,video/x-msvideo
```

### Makefile 追加コマンド

```makefile
.PHONY: shell-minio
shell-minio:
	docker exec -it auction-minio sh

.PHONY: minio-console
minio-console:
	@echo "MinIO Console: http://localhost:9001"
	@echo "Username: minioadmin"
	@echo "Password: minioadmin"

.PHONY: clean-minio
clean-minio:
	docker volume rm real-time-auction_minio_data
```

---

## 3. Go実装

### 必要なライブラリ

```
github.com/minio/minio-go/v7 v7.0.97
github.com/disintegration/imaging v1.6.2  # JPEG処理用（Pure Go実装）
github.com/google/uuid v1.6.0
github.com/aws/aws-sdk-go-v2 v1.21.0      # S3用（本番環境）
github.com/aws/aws-sdk-go-v2/service/s3 v1.40.0
cloud.google.com/go/storage v1.33.0       # GCS用（本番環境）
```

**削除したライブラリ:**
- ~~`github.com/chai2010/webp v1.1.1`~~ - CGO依存のため削除
- ~~`github.com/kolesa-team/go-webp v1.0.4`~~ - CGO依存のため削除

### ディレクトリ構造

```
backend/internal/
├── domain/
│   └── item_media.go              # ItemMedia構造体
├── repository/
│   └── item_media_repository.go   # DB操作
├── service/
│   ├── storage_service.go         # インターフェース
│   ├── image_processor.go         # 画像処理
│   └── storage/
│       ├── minio.go               # MinIO実装
│       ├── s3.go                  # S3実装
│       └── gcs.go                 # GCS実装
├── handler/
│   └── media_handler.go           # エンドポイント
└── middleware/
    └── upload.go                  # ファイルサイズ検証

backend/pkg/utils/
├── file.go                        # ファイル名サニタイズ
└── mime.go                        # MIMEタイプ検証
```

### 主要コンポーネント

#### Domain層（item_media.go）
- ItemMedia構造体（GORM モデル）
- バリデーション（MediaType、DisplayOrder）

#### Repository層（item_media_repository.go）
- Create, FindByItemID, FindByID, Delete
- UpdateDisplayOrder（一括更新）
- CountByItemIDAndType（枚数制限チェック）

#### Service層

**StorageService インターフェース:**
- Upload(ctx, bucket, objectName, filePath, contentType) → (url, error)
- Delete(ctx, bucket, objectName) → error
- GetPublicURL(bucket, objectName) → string
- HealthCheck(ctx) → error

**ImageProcessor:**
- ProcessImage(srcPath) → JPEG変換+リサイズ（品質80、最大1920x1080）
- GenerateThumbnail(srcPath) → 300x300サムネイル（正方形、中央クロップ）
- CleanupTempFiles(paths) → 一時ファイル削除
- ValidateImageFile(filePath) → 画像ファイルバリデーション

#### Handler層（media_handler.go）

**エンドポイント:**

| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| POST | `/api/admin/items/:id/media` | アップロード | auctioneer/admin |
| DELETE | `/api/admin/items/:id/media/:mediaId` | 削除 | auctioneer/admin |
| PUT | `/api/admin/items/:id/media/reorder` | 順序変更 | auctioneer/admin |
| GET | `/api/items/:id/media` | 一覧取得 | 不要 |

---

## 4. API仕様

### POST /api/admin/items/:id/media

**リクエスト:**
```
Content-Type: multipart/form-data
Authorization: Bearer <JWT>

file: (binary)
media_type: "image" or "video"
```

**成功レスポンス（201）:**
```json
{
  "id": 123,
  "item_id": 45,
  "media_type": "image",
  "url": "http://localhost:9000/auction-media/items/45/original_abc123.jpg",
  "thumbnail_url": "http://localhost:9000/auction-media/items/45/thumb_abc123.jpg",
  "display_order": 1,
  "created_at": "2025-12-01T10:30:00Z"
}
```

**エラーレスポンス:**

| ステータス | エラーコード | メッセージ |
|-----------|------------|----------|
| 400 | invalid_file | ファイルが選択されていません |
| 400 | file_too_large | ファイルサイズが上限を超えています |
| 400 | media_limit_exceeded | メディアの登録上限に達しています（画像: 10枚） |
| 403 | forbidden | 権限がありません |
| 500 | processing_failed | 画像処理に失敗しました |

---

## 5. セキュリティ

### 認証・認可
- 全管理系エンドポイントでJWT必須
- ロール: auctioneer, system_admin のみアップロード可能

### ファイルバリデーション
- MIMEタイプ検証（マジックバイト）
- ファイル名サニタイズ（UUID使用）
- パストラバーサル対策
- ファイルサイズ制限（middleware）

### ストレージアクセス制御
- MinIO: 読み取りパブリック、書き込みAPI経由のみ
- S3/GCS: IAMロールで制限

---

## 6. エラーハンドリング

### ロールバックフロー

```
画像処理失敗 → 一時ファイル削除 → エラーレスポンス
  ↓
ストレージアップロード失敗 → 一時ファイル削除 → エラーレスポンス
  ↓
DB保存失敗 → ストレージ削除 → 一時ファイル削除 → エラーレスポンス
```

### 孤立ファイル対策
- DB削除失敗時はストレージ削除しない
- 定期的なクリーンアップジョブ（将来実装）

---

## 7. 実装手順（優先度順）

### Phase 1: 環境構築（2時間）
- [x] docker-compose.ymlにMinIO追加
- [x] .env設定追加
- [x] Makefile更新
- [x] `make up`で起動確認
- [x] MinIO Console（http://localhost:9001）アクセス確認

### Phase 2: Domain・Repository層（3時間）
- [x] internal/domain/item_media.go 実装
- [x] internal/repository/item_media_repository.go 実装
- [x] 単体テスト

### Phase 3: Storage Service（4時間）
- [x] internal/service/storage_service.go インターフェース定義
- [x] internal/service/storage/minio.go 実装
- [x] internal/service/storage/s3.go 骨格実装
- [x] MinIO接続テスト

### Phase 4: 画像処理Service（5時間）
- [x] internal/service/image_processor.go 実装
- [x] リサイズ、JPEG変換、サムネイル生成（WebP→JPEG変更）
- [x] 画像処理テスト（全8テスト成功）
- [x] CGO依存の問題を解決（Pure Go実装に変更）

### Phase 5: Middleware（2時間）
- [x] internal/middleware/upload.go 実装
- [x] pkg/utils/file.go（ファイル名サニタイズ）
- [x] pkg/utils/mime.go（MIMEタイプ検証）
- [x] ユニットテスト作成（全テスト成功）

### Phase 6: Handler実装（6時間）
- [x] internal/handler/media_handler.go 実装
- [x] ルーティング設定（cmd/api/main.go）
- [x] エラーハンドリング
- [x] ロールバック処理

### Phase 7: 統合テスト（4時間）
- [ ] Postmanでエンドポイントテスト
- [ ] エラーケース確認
- [ ] ロールバック動作確認

### Phase 8: フロントエンド対応（8時間）
- [ ] frontend/src/services/media.service.ts
- [ ] frontend/src/components/MediaUploader.vue
- [ ] frontend/src/components/MediaGallery.vue
- [ ] アイテム登録画面への統合

### Phase 9: 本番環境対応（3時間）
- [ ] S3実装完成
- [ ] GCS実装完成
- [ ] 環境変数テンプレート作成

### Phase 10: ドキュメント整備（2時間）
- [ ] API仕様書更新
- [ ] README更新
- [ ] トラブルシューティングガイド

**総推定時間: 39時間（約5営業日）**

---

## 8. 成功基準

### 機能要件
- [ ] 画像を5MB以内でアップロード可能
- [x] 自動JPEG変換・リサイズ（1920x1080以内、品質80）
- [x] サムネイル（300x300）自動生成（正方形、中央クロップ）
- [ ] 画像10枚、動画3本まで登録可能
- [ ] 表示順序変更可能
- [ ] 削除でDB・ストレージ両方から削除
- [ ] 環境変数でMinIO/S3切り替え可能

### 非機能要件
- [ ] アップロード処理3秒以内（5MB画像）
- [ ] エラー時の適切なロールバック
- [ ] 同時アップロードの競合制御

### セキュリティ要件
- [ ] JWT認証必須
- [ ] MIMEタイプ検証（マジックバイト）
- [ ] ファイル名サニタイズ
- [ ] パストラバーサル防止

---

## 9. 重要ファイル

実装時に変更・作成が必要な主要ファイル:

### インフラ
- `docker-compose.yml` - MinIOコンテナ追加
- `.env` - ストレージ設定追加
- `Makefile` - MinIO関連コマンド追加
- `nginx/nginx.conf` - メディア配信ルート追加（必要に応じて）

### バックエンド（新規作成）
- `backend/internal/domain/item_media.go`
- `backend/internal/repository/item_media_repository.go`
- `backend/internal/service/storage_service.go`
- `backend/internal/service/image_processor.go`
- `backend/internal/service/storage/minio.go`
- `backend/internal/service/storage/s3.go`
- `backend/internal/handler/media_handler.go`
- `backend/internal/middleware/upload.go`
- `backend/pkg/utils/file.go`
- `backend/pkg/utils/mime.go`

### バックエンド（既存修正）
- `backend/cmd/api/main.go` - ルーティング追加
- `backend/go.mod` - 依存ライブラリ追加

### フロントエンド（新規作成）
- `frontend/src/services/media.service.ts`
- `frontend/src/components/MediaUploader.vue`
- `frontend/src/components/MediaGallery.vue`

### 既存データベース
- マイグレーションファイル変更不要（`item_media`テーブルは実装済み）

---

## 10. 参考情報

### MinIOについて
- **公式サイト**: https://min.io/
- **GitHub**: https://github.com/minio/minio
- **ライセンス**: AGPLv3（完全無料）
- **2025年の変更**: Web UI管理機能は有料化、CLI/APIは完全無料
- **開発環境での使用**: 問題なし（S3互換APIがフル機能で利用可能）

### 参考ドキュメント
- `backend/migrations/001_create_initial_tables.up.sql` - item_mediaテーブル定義
- `docs/database_definition.md` - メディア管理設計方針
- `docs/screen_list.md` - UI要件（画像10枚、動画3本）
- `docs/architecture.md` - 外部ストレージ参照

---

## 11. 技術的な補足

### WebP vs JPEG の比較

| 項目 | WebP（当初計画） | JPEG（実装版） |
|------|----------------|---------------|
| **ファイルサイズ** | 小さい（JPEGの約30%削減） | やや大きい |
| **品質** | 高品質 | 高品質（品質80で十分） |
| **ブラウザ対応** | モダンブラウザのみ | 全ブラウザ対応 |
| **Go実装** | CGO必須（libwebp依存） | Pure Go（imaging） |
| **Dockerビルド** | 複雑（gcc、libwebp-dev必須） | シンプル |
| **クロスコンパイル** | 困難 | 容易 |
| **保守性** | 低い（C依存関係） | 高い（Goのみ） |

### CGO問題の詳細

**発生したエラー:**
```
undefined: webpGetInfo
undefined: webpDecodeGray
undefined: webpEncodeRGB
```

**原因:**
- WebPライブラリ（chai2010/webp、kolesa-team/go-webp）はC言語のlibwebpをCGO経由で呼び出す
- Docker環境でのビルドには以下が必要:
  - Cコンパイラ（gcc、musl-dev）
  - libwebp開発ヘッダー（libwebp-dev）
  - CGOの有効化（`CGO_ENABLED=1`）

**解決策の選択肢:**
1. **Dockerfileに依存関係を追加**（複雑、ビルド時間増加）
   ```dockerfile
   RUN apk add --no-cache gcc musl-dev libwebp-dev
   ```
2. **Pure Go実装に変更**（採用）
   - `github.com/disintegration/imaging` を使用
   - JPEG形式で品質80を設定
   - シンプルで保守性が高い

### 実装済みファイル

**Phase 4完了ファイル:**
- `backend/internal/service/image_processor.go` - 画像処理サービス（168行）
- `backend/internal/service/image_processor_test.go` - テストコード（全8テスト成功）
- `backend/go.mod` - 依存ライブラリ追加（disintegration/imaging v1.6.2）

**テスト結果:**
```
✓ TestNewImageProcessor - インスタンス生成
✓ TestProcessImage - 画像処理全体（リサイズ＋サムネイル）
✓ TestResizeImage - リサイズ機能（4パターン）
✓ TestGenerateThumbnail - サムネイル生成
✓ TestSaveAsJPEG - JPEG保存
✓ TestCleanupTempFiles - 一時ファイル削除
✓ TestValidateImageFile - バリデーション（3パターン）
✓ TestProcessImage_InvalidFile - エラーハンドリング

PASS: 全8テスト (0.195秒)
```
