# オークション編集機能（1.3.3）実装計画

## 概要

既存のオークション作成機能（`AuctionNewView.vue`）をベースに、オークション・商品情報の編集機能を追加。`pending`ステータスのオークションのみ編集可能とし、バックエンドAPIとフロントエンド画面を新規実装する。

**パス**: `/admin/auctions/:id/edit`  
**権限**: auctioneer / system_admin

---

## 実装ステップ

### Step 1: バックエンドDomain層の拡張

**対象ファイル**: `backend/internal/domain/`

**実装内容**:
- `UpdateAuctionRequest` 構造体を追加
- `UpdateItemRequest` 構造体を追加
- `AddItemRequest` 構造体を追加
- バリデーションタグを定義

**構造体定義**:
```go
// UpdateAuctionRequest - オークション更新リクエスト
type UpdateAuctionRequest struct {
    Title       *string `json:"title" validate:"omitempty,max=200"`
    Description *string `json:"description"`
}

// UpdateItemRequest - 商品更新リクエスト
type UpdateItemRequest struct {
    Name        *string `json:"name" validate:"omitempty,max=200"`
    Description *string `json:"description"`
    Metadata    *JSONB  `json:"metadata"`
}

// AddItemRequest - 商品追加リクエスト
type AddItemRequest struct {
    Name        string `json:"name" validate:"required,max=200"`
    Description string `json:"description"`
    Metadata    JSONB  `json:"metadata"`
}
```

---

### Step 2: Repository層に更新メソッド追加

**対象ファイル**: `backend/internal/repository/postgres/`

**実装内容**:
- `UpdateAuction(ctx, id, req)` - オークション更新
- `UpdateItem(ctx, id, req)` - 商品更新
- `DeleteItem(ctx, id)` - 商品削除
- `AddItem(ctx, auctionID, req)` - 商品追加
- `ReorderItems(ctx, auctionID, itemIDs)` - 商品順序変更
- トランザクション対応

**インターフェース定義**:
```go
type AuctionRepository interface {
    // 既存メソッド...
    UpdateAuction(ctx context.Context, id uuid.UUID, req *domain.UpdateAuctionRequest) (*domain.Auction, error)
    GetAuctionWithItems(ctx context.Context, id uuid.UUID) (*domain.AuctionWithItems, error)
}

type ItemRepository interface {
    // 既存メソッド...
    UpdateItem(ctx context.Context, id uuid.UUID, req *domain.UpdateItemRequest) (*domain.Item, error)
    DeleteItem(ctx context.Context, id uuid.UUID) error
    AddItem(ctx context.Context, auctionID uuid.UUID, req *domain.AddItemRequest) (*domain.Item, error)
    ReorderItems(ctx context.Context, auctionID uuid.UUID, itemIDs []uuid.UUID) error
}
```

---

### Step 3: Service層にビジネスロジック実装

**対象ファイル**: `backend/internal/service/`

**実装内容**:
- 編集可否判定（オークションステータス確認）
- 商品削除可否判定（入札状況確認）
- 更新処理（トランザクション）
- `lot_number` 再計算ロジック

**ビジネスルール**:
| 条件 | 編集可否 | 備考 |
|------|---------|------|
| オークション `pending` | ○ | 全項目編集可能 |
| オークション `active` | △ | 説明のみ編集可能、商品追加・削除不可 |
| オークション `ended`/`cancelled` | × | 閲覧のみ |
| 商品 `started_at` が NULL | ○ | 全項目編集可能 |
| 商品 `started_at` 設定済み | × | 編集不可 |
| 商品に入札あり | × | 削除不可 |

**サービスメソッド**:
```go
type AuctionService interface {
    // 既存メソッド...
    GetAuctionForEdit(ctx context.Context, id uuid.UUID) (*AuctionEditResponse, error)
    UpdateAuction(ctx context.Context, id uuid.UUID, req *UpdateAuctionRequest) (*Auction, error)
    CanEditAuction(ctx context.Context, id uuid.UUID) (bool, string, error)
}

type ItemService interface {
    // 既存メソッド...
    UpdateItem(ctx context.Context, id uuid.UUID, req *UpdateItemRequest) (*Item, error)
    DeleteItem(ctx context.Context, id uuid.UUID) error
    AddItem(ctx context.Context, auctionID uuid.UUID, req *AddItemRequest) (*Item, error)
    ReorderItems(ctx context.Context, auctionID uuid.UUID, itemIDs []uuid.UUID) error
    CanEditItem(ctx context.Context, id uuid.UUID) (bool, string, error)
    CanDeleteItem(ctx context.Context, id uuid.UUID) (bool, string, error)
}
```

---

### Step 4: Handler層にエンドポイント追加

**対象ファイル**: `backend/internal/handler/`

**新規エンドポイント**:
| メソッド | パス | 説明 |
|---------|------|------|
| `GET` | `/admin/auctions/:id` | 管理者用オークション詳細取得（編集用） |
| `PUT` | `/admin/auctions/:id` | オークション更新 |
| `PUT` | `/admin/auctions/:id/items/:item_id` | 商品更新 |
| `DELETE` | `/admin/auctions/:id/items/:item_id` | 商品削除 |
| `POST` | `/admin/auctions/:id/items` | 商品追加 |
| `PUT` | `/admin/auctions/:id/items/reorder` | 商品順序変更 |

**レスポンス形式**:
```json
// GET /admin/auctions/:id
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "pending|active|ended|cancelled",
  "started_at": "datetime|null",
  "can_edit": true,
  "can_edit_reason": "string|null",
  "items": [
    {
      "id": "uuid",
      "name": "string",
      "description": "string",
      "metadata": {},
      "lot_number": 1,
      "starting_price": null,
      "current_price": null,
      "started_at": null,
      "ended_at": null,
      "can_edit": true,
      "can_delete": true,
      "bid_count": 0
    }
  ],
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

---

### Step 5: フロントエンドAPI・Store拡張

**対象ファイル**: 
- `frontend/src/api/auction.js`
- `frontend/src/stores/auction.js`

**API追加**:
```javascript
// auction.js
export const getAuctionForEdit = (id) => api.get(`/admin/auctions/${id}`)
export const updateAuction = (id, data) => api.put(`/admin/auctions/${id}`, data)
export const updateItem = (auctionId, itemId, data) => api.put(`/admin/auctions/${auctionId}/items/${itemId}`, data)
export const deleteItem = (auctionId, itemId) => api.delete(`/admin/auctions/${auctionId}/items/${itemId}`)
export const addItem = (auctionId, data) => api.post(`/admin/auctions/${auctionId}/items`, data)
export const reorderItems = (auctionId, itemIds) => api.put(`/admin/auctions/${auctionId}/items/reorder`, { item_ids: itemIds })
```

**Store追加**:
```javascript
// stores/auction.js
const currentAuction = ref(null)
const isLoadingAuction = ref(false)

async function fetchAuctionForEdit(id) { ... }
async function handleUpdateAuction(id, data) { ... }
async function handleUpdateItem(auctionId, itemId, data) { ... }
async function handleDeleteItem(auctionId, itemId) { ... }
async function handleAddItem(auctionId, data) { ... }
async function handleReorderItems(auctionId, itemIds) { ... }
```

---

### Step 6: 編集画面コンポーネント作成

**対象ファイル**: `frontend/src/views/admin/AuctionEditView.vue`

**コンポーネント構成**:
```
AuctionEditView.vue
├── AuctionBasicInfoForm.vue（既存を再利用）
├── AuctionItemList.vue（既存を拡張）
│   └── AuctionItemForm.vue（既存を再利用）
├── AuctionStatusBadge.vue（新規）
└── DeleteConfirmModal.vue（新規）
```

**主要機能**:
1. **データ取得・表示**
   - オークションIDからデータ取得
   - ローディング状態表示
   - エラー表示（404等）

2. **編集可否判定**
   - `can_edit` フラグに基づくフォーム制御
   - 編集不可時はメッセージ表示

3. **オークション基本情報編集**
   - タイトル、説明の編集
   - onBlur時バリデーション

4. **商品管理**
   - 商品の追加・編集・削除
   - 順序変更（▲▼ボタン）
   - 削除確認モーダル

5. **更新処理**
   - 変更差分検出
   - 一括更新 or 個別更新

**状態管理**:
```javascript
const auctionId = ref(route.params.id)
const originalData = ref(null)  // 変更検出用
const formData = ref({ title: '', description: '', items: [] })
const errors = ref({ title: '', description: '' })
const itemErrors = ref([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const canEdit = ref(false)
const canEditReason = ref('')
const hasChanges = computed(() => { /* 差分検出 */ })
```

---

### Step 7: ルーティング設定追加

**対象ファイル**: `frontend/src/router/index.js`

**追加ルート**:
```javascript
{
  path: '/admin/auctions/:id/edit',
  name: 'auction-edit',
  component: () => import('../views/admin/AuctionEditView.vue'),
  meta: { 
    requiresAuth: true, 
    requireAdminOrAuctioneer: true,
    title: 'オークション編集'
  }
}
```

---

## 画面レイアウト（ワイヤーフレーム）

```
┌─────────────────────────────────────────────────────────────┐
│ ← 戻る                    オークション編集                    │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ ステータス: [pending] 編集可能                            │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─ 基本情報 ─────────────────────────────────────────────┐ │
│ │ タイトル *                                              │ │
│ │ ┌─────────────────────────────────────────────────────┐│ │
│ │ │ オークション名                                       ││ │
│ │ └─────────────────────────────────────────────────────┘│ │
│ │                                                        │ │
│ │ 説明                                                   │ │
│ │ ┌─────────────────────────────────────────────────────┐│ │
│ │ │                                                     ││ │
│ │ │                                                     ││ │
│ │ └─────────────────────────────────────────────────────┘│ │
│ └────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─ 商品一覧 ─────────────────────────────────────────────┐ │
│ │ [+ 商品を追加]                                          │ │
│ │                                                        │ │
│ │ ┌─ 商品 1 ──────────────────────────────────────────┐ │ │
│ │ │ ▲ ▼  商品名: [競走馬A        ]  [編集] [削除]     │ │ │
│ │ │      説明: サラブレッド、3歳牡馬...                │ │ │
│ │ └──────────────────────────────────────────────────┘ │ │
│ │                                                        │ │
│ │ ┌─ 商品 2 ──────────────────────────────────────────┐ │ │
│ │ │ ▲ ▼  商品名: [競走馬B        ]  [編集] [削除]     │ │ │
│ │ │      説明: サラブレッド、4歳牝馬...                │ │ │
│ │ └──────────────────────────────────────────────────┘ │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │              [キャンセル]  [更新する]                    │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## バリデーションルール

### オークション基本情報

| フィールド | 必須 | ルール |
|-----------|------|--------|
| title | ○ | 1〜200文字 |
| description | × | 制限なし |

### 商品情報

| フィールド | 必須 | ルール |
|-----------|------|--------|
| name | ○ | 1〜200文字 |
| description | × | 制限なし |
| metadata | × | 有効なJSON |

---

## エラーハンドリング

| エラーケース | ステータス | メッセージ |
|-------------|-----------|-----------|
| オークションが見つからない | 404 | オークションが見つかりません |
| 編集権限がない | 403 | このオークションを編集する権限がありません |
| オークションが編集不可状態 | 400 | 開始後のオークションは編集できません |
| 商品が削除不可（入札あり） | 400 | 入札のある商品は削除できません |
| バリデーションエラー | 422 | 入力内容に誤りがあります |
| 競合エラー（楽観的ロック） | 409 | 他のユーザーが更新しました。再読み込みしてください |

---

## テストケース

### バックエンドユニットテスト

1. **オークション更新**
   - `pending` ステータスで更新成功
   - `active` ステータスで更新失敗
   - 存在しないIDで404

2. **商品更新**
   - 未開始商品の更新成功
   - 開始済み商品の更新失敗

3. **商品削除**
   - 入札なし商品の削除成功
   - 入札あり商品の削除失敗

4. **商品追加**
   - `pending` オークションへの追加成功
   - `active` オークションへの追加失敗
   - `lot_number` 自動採番確認

### フロントエンドテスト

1. **表示テスト**
   - データ取得・表示
   - 編集不可時のUI表示

2. **フォームテスト**
   - バリデーション
   - 変更検出

3. **操作テスト**
   - 商品追加・削除・順序変更
   - 更新処理

---

## 工数見積もり

| 項目 | 工数 |
|-----|------|
| Step 1: Domain層拡張 | 1時間 |
| Step 2: Repository層実装 | 2時間 |
| Step 3: Service層実装 | 2時間 |
| Step 4: Handler層実装 | 2時間 |
| Step 5: フロントエンドAPI・Store | 2時間 |
| Step 6: 編集画面コンポーネント | 4時間 |
| Step 7: ルーティング設定 | 0.5時間 |
| テスト実装 | 3時間 |
| **合計** | **16.5時間** |

---

## 考慮事項

### 楽観的ロック（将来実装）

複数管理者が同時編集した場合の競合回避:
- 初期実装: `updated_at` での簡易的な検証
- 将来: `version` カラム追加による厳密なロック

### 商品順序の再計算

商品追加・削除時に `lot_number` を再計算:
- 追加時: 最大 `lot_number + 1`
- 削除後: 再採番 or 欠番許容

### メディア管理（将来実装）

現在はメディアアップロード機能未実装のため、編集画面でも未対応:
- 将来: ドラッグ&ドロップでのアップロード
- 将来: 画像・動画の順序変更

---

## 関連ドキュメント

- [画面一覧](./screen_list.md) - 1.3.3 オークション編集
- [データベース定義](./database_definition.md) - auctions, items テーブル
- [オークション作成実装計画](./plan/auction_creation_implementation_plan.md)
