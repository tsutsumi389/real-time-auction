# Plan: 管理者向け商品管理画面の作成（最終版）

商品（Item）をオークションから独立して管理できる画面を追加します。商品一覧・新規作成・編集画面と、オークション商品紐づけ画面を実装します。商品はオークション未割当の状態でも登録可能で、紐づけ画面からオークションへの追加・解除（開始前のみ）を行います。

## Steps

### Step 1: Backend - DBマイグレーション追加

**ファイル**: `backend/migrations/014_make_item_auction_nullable.up.sql`（新規作成）

`items.auction_id` を `NULL` 許可に変更し、`lot_number` のユニーク制約を `(auction_id, lot_number) WHERE auction_id IS NOT NULL` の部分インデックスに修正。

---

### Step 2: Backend - Domainモデル更新

**ファイル**: `backend/internal/domain/item.go`

- `AuctionID` を `*uuid.UUID`（ポインタ型）に変更
- `CreateItemRequest`（オークション無し作成用）型を追加
- `AssignItemsRequest` 型を追加
- `UnassignItemRequest` 型を追加

---

### Step 3: Backend - 商品管理API追加

**ファイル**: `backend/internal/handler/item_handler.go`（新規作成）

以下のエンドポイントを実装:

| メソッド | パス | 機能 |
|---------|------|------|
| GET | `/api/admin/items` | 一覧（フィルタ: `assigned`/`unassigned`/`all`、キーワード検索、ページネーション） |
| GET | `/api/admin/items/:id` | 詳細 |
| POST | `/api/admin/items` | 新規作成（`auction_id = NULL`） |
| PUT | `/api/admin/items/:id` | 更新 |
| DELETE | `/api/admin/items/:id` | 削除（オークション紐づき or 入札履歴ありは403エラー） |

**権限**: `system_admin` / `auctioneer`

---

### Step 4: Backend - 紐づけAPI追加

**ファイル**: `backend/internal/handler/auction_handler.go`（既存に追加）

| メソッド | パス | 機能 |
|---------|------|------|
| POST | `/api/admin/auctions/:id/items/assign` | 商品を追加（商品ID配列、`lot_number` 自動採番） |
| DELETE | `/api/admin/auctions/:id/items/:itemId/unassign` | 商品を解除（オークション開始後は403エラー、`auction_id = NULL`, `lot_number = 0`） |

**ビジネスルール**:
- 紐づけ時: 未割当商品のみ追加可能、`lot_number` は既存最大値+1から自動採番
- 解除時: オークション開始前（全商品が `started_at = NULL`）のみ可能

---

### Step 5: Frontend - 商品管理画面作成

**ファイル**:
- `frontend/src/views/admin/ItemListView.vue`（新規）
- `frontend/src/views/admin/ItemNewView.vue`（新規）
- `frontend/src/views/admin/ItemEditView.vue`（新規）

**商品一覧画面** (`/admin/items`):
- フィルタ: 未割当 / 割当済み / すべて
- キーワード検索（商品名）
- ページネーション（20件/ページ）
- テーブル列: ID（短縮）、商品名、紐づけオークション名（リンク）、開始価格、作成日
- アクション: 編集ボタン

**商品新規作成画面** (`/admin/items/new`):
- 入力項目: 商品名（必須）、説明（任意）、開始価格（任意）
- 保存ボタン → 一覧に戻る

**商品編集画面** (`/admin/items/:id/edit`):
- 入力項目: 商品名、説明、開始価格
- 紐づけオークション表示（リンク付き）、未割当の場合は「未割当」表示
- 削除ボタン（紐づきありまたは入札履歴ありの場合は無効化 + ツールチップで理由表示）
- 保存ボタン → 一覧に戻る

---

### Step 6: Frontend - オークション商品紐づけ画面作成

**ファイル**: `frontend/src/views/admin/AuctionItemsAssignView.vue`（新規）

**パス**: `/admin/auctions/:id/items`

**レイアウト**: 2カラム

**左カラム - 未割当商品一覧**:
- 検索ボックス（商品名）
- 商品リスト（チェックボックス付き）
- 「選択した商品を追加」ボタン

**右カラム - オークション内商品一覧**:
- 商品リスト（ロット番号順）
- ドラッグ&ドロップで順序変更（`lot_number` 更新）
- 各商品に「解除」ボタン
- オークション開始後は解除ボタン無効化

**遷移元**: オークション編集画面から「商品管理」ボタン

---

### Step 7: Frontend - Store・API・Router設定

**新規ファイル**:
- `frontend/src/stores/item.js` - Pinia Store
- `frontend/src/services/itemApi.js` - API サービス

**既存ファイル変更**:
- `frontend/src/services/auctionApi.js` - 紐づけAPI追加
- `frontend/src/router/index.js` - ルート追加

**ルート追加**:
- `/admin/items` → `ItemListView.vue`
- `/admin/items/new` → `ItemNewView.vue`
- `/admin/items/:id/edit` → `ItemEditView.vue`
- `/admin/auctions/:id/items` → `AuctionItemsAssignView.vue`

**ナビゲーション更新**:
- 管理者メニューに「商品管理」を追加（オークション一覧の下）

---

## ビジネスルール詳細

### 商品削除の制約
- オークションに紐づいている商品 → 削除不可（403エラー）
- 入札履歴がある商品 → 削除不可（403エラー）
- 上記以外 → 削除可能

### オークション商品解除の制約
- オークション内のいずれかの商品が開始済み（`started_at IS NOT NULL`）→ 解除不可（403エラー）
- 全商品が未開始 → 解除可能（`auction_id = NULL`, `lot_number = 0` に更新）

### 商品編集画面の表示
- 紐づけオークション: オークション名をリンク表示（`/admin/auctions/:id/edit` へ遷移可能）
- 未割当の場合: 「未割当」と表示

---

## API仕様

### GET /api/admin/items

**Query Parameters**:
- `status`: `assigned` | `unassigned` | `all` (default: `all`)
- `search`: キーワード検索
- `page`: ページ番号 (default: 1)
- `limit`: 件数 (default: 20)

**Response**:
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "商品名",
      "description": "説明",
      "starting_price": 100000,
      "auction_id": "uuid or null",
      "auction_title": "オークション名 or null",
      "lot_number": 1,
      "bid_count": 0,
      "can_delete": true,
      "created_at": "2025-11-30T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

### POST /api/admin/items

**Request**:
```json
{
  "name": "商品名",
  "description": "説明",
  "starting_price": 100000
}
```

### POST /api/admin/auctions/:id/items/assign

**Request**:
```json
{
  "item_ids": ["uuid1", "uuid2"]
}
```

### DELETE /api/admin/auctions/:id/items/:itemId/unassign

**Response**: 204 No Content

---

## 実装順序

1. Backend: マイグレーション作成・適用
2. Backend: Domain更新
3. Backend: Handler/Service/Repository 実装
4. Frontend: Store/API作成
5. Frontend: 商品一覧画面
6. Frontend: 商品新規・編集画面
7. Frontend: オークション商品紐づけ画面
8. Frontend: Router・ナビゲーション更新
9. テスト・動作確認
