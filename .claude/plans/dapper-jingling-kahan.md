# オークション管理UI/UX改善実装プラン

## 概要

オークション作成時の商品登録を不要にし、作成後に商品を紐付ける柔軟なワークフローへ変更します。商品割り当てインターフェースと編集画面を改善し、操作性を向上させます。

## 主な変更内容

1. **オークション作成の簡素化**: 商品なしで作成可能に
2. **作成後の編集画面へリダイレクト**: 商品追加への明確な導線
3. **編集画面のタブ化**: 新規作成・既存選択を統一インターフェースで提供
4. **ドラッグ&ドロップの改善**: 視覚フィードバックとモバイル対応強化

## 実装フェーズ

### Phase 1: バックエンド変更（必須・最優先）

#### 1.1 ドメイン層のバリデーション変更

**ファイル**: `backend/internal/domain/auction.go`

**変更箇所**: 77行目
```go
// 変更前
Items []CreateItemRequest `json:"items" binding:"required,min=1,dive"`

// 変更後
Items []CreateItemRequest `json:"items" binding:"omitempty,dive"`
```

**理由**:
- `omitempty`: 空の商品配列を許可（0件でもOK）
- `dive`: 商品が提供された場合は各要素をバリデーション
- 後方互換性: 既存の商品付きリクエストも引き続き動作

#### 1.2 サービス層の更新

**ファイル**: `backend/internal/service/auction_service.go`

**関数**: `CreateAuction` (259-298行目)

**必要な変更**:
1. 空の商品配列を適切に処理
2. `len(req.Items) > 0` の場合のみ商品を作成
3. トランザクションの整合性を維持
4. レスポンスで `item_count = 0` を返す

#### 1.3 リポジトリ層の更新

**ファイル**: `backend/internal/repository/auction_repository.go`

**関数**: `CreateAuctionWithItems`

**必要な変更**:
- 商品配列が空の場合、商品の挿入をスキップ
- オークションレコードは作成
- 成功時は nil エラーを返す

#### 1.4 テストの追加

**ファイル**: `backend/internal/service/auction_service_test.go`

**テストケース**:
- ✅ 0件の商品でオークション作成
- ✅ 1件の商品でオークション作成（既存の動作確認）
- ✅ 複数の商品でオークション作成（既存の動作確認）
- ✅ バリデーションエラー: 空のタイトル
- ✅ バリデーションエラー: タイトル長すぎ（>200文字）

---

### Phase 2: オークション作成画面の簡素化

#### 2.1 AuctionNewView.vueの更新

**ファイル**: `frontend/src/views/admin/AuctionNewView.vue`

**主な変更**:

1. **ItemListコンポーネントの削除**:
   - インポート削除（8行目付近）
   - テンプレートから削除（217-222行目）
   - items ref と関連バリデーションロジック削除

2. **テンプレート更新**:
   - 商品セクションを削除
   - 情報メッセージを追加: 「商品はオークション作成後に追加できます」
   - AuctionBasicInfoコンポーネントのみ保持

3. **送信ハンドラの更新**:
   - リダイレクト先を変更: `auction-list` → `auction-edit/:id`
   - クエリパラメータ追加: `?created=true`
   ```javascript
   router.push({
     name: 'auction-edit',
     params: { id: result.id },
     query: { created: 'true' }
   })
   ```

4. **バリデーション更新**:
   - 商品バリデーション削除（119-130行目）
   - title, description, started_at のみバリデーション

#### 2.2 ストアの確認

**ファイル**: `frontend/src/stores/auction.js`

**確認事項**: `handleCreateAuction` が完全なレスポンスオブジェクト（id含む）を返すことを確認
- 既存実装で問題なし（254行目で `return response`）

---

### Phase 3: 編集画面の成功メッセージとCTA

#### 3.1 AuctionEditView.vueの更新

**ファイル**: `frontend/src/views/admin/AuctionEditView.vue`

**追加機能**:

1. **クエリパラメータ検出** (onMounted内、78行目付近):
   ```javascript
   if (route.query.created === 'true') {
     showSuccessMessage.value = true
     defaultTab.value = 'create'
   }
   ```

2. **成功メッセージバナー** (テンプレート冒頭):
   ```vue
   <div v-if="showSuccessMessage" class="success-banner">
     ✅ オークションを作成しました
     次に商品を追加してオークションを完成させましょう
     <Button @click="showSuccessMessage = false">閉じる</Button>
   </div>
   ```

3. **自動消去**: 5秒後に自動的にメッセージを消去

---

### Phase 4: 編集画面のタブインターフェース（優先度：高）

#### 4.1 タブナビゲーションの追加

**ファイル**: `frontend/src/views/admin/AuctionEditView.vue`

**使用コンポーネント**: Radix Vue Tabs

**タブ構成**:
1. **割当済み商品** (`assigned`): 既存のEditItemListコンポーネント
2. **新規作成** (`create`): 新規作成用のインラインフォーム
3. **既存から選択** (`select`): 未割当商品の選択UI

**追加位置**: 503行目以降のQuick Actionカードの後

#### 4.2 新規コンポーネント: ItemCreateInline.vue

**ファイル**: `frontend/src/components/admin/ItemCreateInline.vue`（新規作成）

**機能**:
- 商品情報入力フォーム（名前・説明・開始価格）
- バリデーション（名前必須）
- 作成 → 割り当ての2段階処理
- ローディング状態表示

**Props**:
- `auctionId` (String, required)
- `loading` (Boolean)

**Emits**:
- `item-created`: 商品作成成功時

**処理フロー**:
1. `POST /api/admin/items` で商品作成（auction_id = null）
2. 作成した商品IDを取得
3. `POST /api/admin/auctions/:id/items/assign` で割り当て
4. 成功後、`item-created` イベントを emit
5. 親コンポーネント（AuctionEditView）が「割当済み商品」タブに切り替え

#### 4.3 新規コンポーネント: ItemSelectList.vue

**ファイル**: `frontend/src/components/admin/ItemSelectList.vue`（新規作成）

**機能**:
- 未割当商品一覧の取得・表示
- 検索・フィルター機能
- チェックボックスで複数選択
- ページネーション
- 一括割り当てボタン

**Props**:
- `auctionId` (String, required)
- `selectedIds` (Array)

**Emits**:
- `update:selectedIds`: 選択状態の更新
- `items-assigned`: 割り当て成功時

**API呼び出し**:
- `GET /api/admin/items?status=unassigned&search=...&page=...`
- `POST /api/admin/auctions/:id/items/assign`

#### 4.4 EditItemList.vueのドラッグ&ドロップ追加

**ファイル**: `frontend/src/components/admin/EditItemList.vue`

**変更内容**:

1. **上下ボタンの置き換え**: ドラッグハンドルアイコン（⋮⋮）に変更
2. **ドラッグイベントハンドラ追加**:
   - `handleDragStart`: ドラッグ開始時にインデックス保存
   - `handleDragOver`: ドロップゾーン表示
   - `handleDrop`: 並び替え実行 + API呼び出し

3. **視覚フィードバック**:
   - ドラッグ中: opacity-50, scale-95
   - ドロップゾーン: 青い点線ボーダー
   - ホバー: 背景色変更
   - トランジション: 200ms

**参考実装**: AuctionItemsAssignView.vue の 468-514行目のパターンを流用

---

### Phase 5: 商品割り当て画面の視覚改善（優先度：最高）

#### 5.1 ドラッグ&ドロップの視覚フィードバック強化

**ファイル**: `frontend/src/views/admin/AuctionItemsAssignView.vue`

**改善内容**:

1. **ドラッグハンドルの明確化** (194-254行目):
   - アイコン: 三本線（⋮⋮）を追加
   - サイズ: w-5 h-5
   - カラー: gray-400
   - カーソル: cursor-move on hover

2. **ドラッグ中の視覚状態**:
   ```css
   .dragging {
     opacity: 0.5;
     transform: scale(0.95);
     box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
     z-index: 50;
   }
   ```

3. **ドロップゾーンインジケーター**:
   ```css
   .drop-zone {
     border: 2px dashed rgb(96, 165, 250); /* blue-400 */
     padding: 8px;
     background: rgb(239, 246, 255); /* blue-50 */
     animation: pulse 1s infinite;
   }
   ```

4. **ホバー状態**:
   ```css
   .item-hover {
     background: rgb(239, 246, 255); /* blue-50 */
     border-color: rgb(191, 219, 254); /* blue-200 */
     transition: all 200ms;
   }
   ```

#### 5.2 ローディング状態の改善

**追加要素**:
- スピナーオーバーレイ（並び替えAPI呼び出し中）
- ドラッグ無効化（ローディング中）
- プログレスインジケーター

#### 5.3 モバイル対応（タッチイベント）

**追加ハンドラ**:
- `handleTouchStart`: タッチ開始 → ドラッグ開始に変換
- `handleTouchMove`: タッチ移動 → ドラッグ位置更新
- `handleTouchEnd`: タッチ終了 → ドロップ実行

**レスポンシブ調整**:
- 768px未満: カラムを縦積み
- ドラッグハンドルサイズ拡大（モバイルは w-8 h-8）

---

### Phase 6: オークション一覧のクイックアクション（優先度：低）

#### 6.1 ドロップダウンメニューの追加

**ファイル**: `frontend/src/components/admin/AuctionTable.vue`

**使用コンポーネント**: Radix Vue DropdownMenu

**メニュー項目**:
- 編集（既存）
- 商品追加（新規）→ `/admin/auctions/:id/edit?tab=create`
- ライブ表示（既存）
- 開始/終了/キャンセル（既存、条件付き）

#### 6.2 商品数バッジの追加

**表示位置**: タイトル列の横

**表示形式**:
- 0件: グレーバッジ「🏷️ 0件」
- 1件以上: ブルーバッジ「🏷️ 5件」

---

## 重要ファイル一覧

### バックエンド
1. `backend/internal/domain/auction.go` - バリデーション変更（77行目）
2. `backend/internal/service/auction_service.go` - サービス層ロジック（259-298行目）
3. `backend/internal/repository/auction_repository.go` - リポジトリ層
4. `backend/internal/handler/auction_handler.go` - ハンドラ（確認のみ）
5. `backend/internal/service/auction_service_test.go` - テスト追加

### フロントエンド
1. `frontend/src/views/admin/AuctionNewView.vue` - 作成画面簡素化
2. `frontend/src/views/admin/AuctionEditView.vue` - タブUI追加、成功メッセージ
3. `frontend/src/views/admin/AuctionItemsAssignView.vue` - 視覚改善
4. `frontend/src/components/admin/EditItemList.vue` - ドラッグ&ドロップ追加
5. `frontend/src/components/admin/ItemCreateInline.vue` - 新規作成（新規ファイル）
6. `frontend/src/components/admin/ItemSelectList.vue` - 選択UI（新規ファイル）
7. `frontend/src/components/admin/AuctionTable.vue` - クイックアクション追加
8. `frontend/src/stores/auction.js` - 確認のみ

---

## API変更

### 変更されるエンドポイント

**POST /api/admin/auctions**

**リクエストボディ（変更後）**:
```json
{
  "title": "秋の馬セリ",
  "description": "秋シーズンの馬セリ開催",
  "started_at": "2025-12-15T10:00:00Z",
  "items": []  ← 空配列OK
}
```

**レスポンス（変更なし）**:
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "title": "秋の馬セリ",
  "status": "pending",
  "item_count": 0,
  "created_at": "2025-12-05T10:00:00Z"
}
```

### 既存エンドポイント（変更なし）

- `POST /api/admin/items` - 未割当商品作成
- `POST /api/admin/auctions/:id/items/assign` - 商品割り当て
- `PUT /api/admin/auctions/:id/items/reorder` - 並び替え
- `GET /api/admin/items?status=unassigned` - 未割当商品一覧

---

## 実装順序（推奨）

### Week 1: コア機能（必須）
1. **Phase 1**: バックエンド変更（2.5時間）
2. **Phase 2**: オークション作成簡素化（3.5時間）
3. **Phase 3**: 編集画面成功メッセージ（1時間）

### Week 2: UI/UX改善（優先度高）
4. **Phase 5**: 商品割り当て画面の視覚改善（4.5時間）
5. **Phase 4**: 編集画面タブインターフェース（7時間）

### Week 3: 仕上げ（優先度低）
6. **Phase 6**: オークション一覧クイックアクション（2.5時間）
7. **テスト・ドキュメント**: 統合テスト、ドキュメント更新（3.5時間）

**合計見積もり時間**: 約24時間

---

## 成功基準

### 機能要件
- ✅ 商品なしでオークションを作成できる
- ✅ 作成後、編集画面にリダイレクトされる
- ✅ 編集画面で成功メッセージと「商品追加」CTAが表示される
- ✅ 編集画面に3つのタブ（割当済み・新規作成・既存選択）がある
- ✅ インライン商品作成がページリロードなしで動作する
- ✅ 既存商品を選択して一括割り当てできる
- ✅ ドラッグ&ドロップでスムーズに並び替えできる
- ✅ モバイルでタッチ操作できる

### UX要件
- ✅ ドラッグ&ドロップが即座に反応（<50ms）
- ✅ API呼び出しが2秒以内に完了
- ✅ ページ遷移がスムーズ（<300ms）
- ✅ 成功/エラーメッセージが明確
- ✅ ローディング状態が可視化されている

### アクセシビリティ
- ✅ キーボードナビゲーションが動作する
- ✅ スクリーンリーダーが変更を通知する
- ✅ フォーカス管理が適切

---

## ロールバック戦略

### バックエンド
問題が発生した場合、`backend/internal/domain/auction.go` 77行目を元に戻す:
```go
Items []CreateItemRequest `json:"items" binding:"required,min=1,dive"`
```
再デプロイするだけでロールバック完了。

### フロントエンド
- Phase 1-2のみデプロイ可能（Phase 3-6は独立）
- Phase 3-5は個別にロールバック可能
- git revert で該当コミットを元に戻す

---

## データベース変更

**変更なし**

既存のスキーマで対応可能:
- `auctions` テーブル: 商品数0でも問題なし
- `items.auction_id`: 既にNULL許容（未割当状態）
- 既存のインデックスと制約は有効なまま

---

## セキュリティ

**新たなリスクなし**:
- 認証・認可は変更なし
- RBACの適用（auctioneer/system_admin のみ）
- CSRFトークン保護（既存ミドルウェア）
- 入力バリデーション（タイトル必須、最大長チェック）
- SQL インジェクション対策（GORM ORM使用）
- XSS対策（Vue自動エスケープ）

---

## 監視とメトリクス

### 追跡すべき指標

**ビジネスメトリクス**:
- 商品0件で作成されたオークションの割合
- 作成後、最初の商品追加までの平均時間
- オークションあたりの平均商品数（変更前後の比較）

**技術メトリクス**:
- API レスポンスタイム（p50, p95, p99）
- フロントエンドエラー率
- ドラッグ&ドロップ成功率
- モバイル vs デスクトップ利用率

### アラート設定

**クリティカル**:
- API エラー率 > 5%
- 平均レスポンスタイム > 2秒
- フロントエンドクラッシュ率 > 1%

---

## 今後の拡張（将来的）

- 一括商品編集（複数の価格を一度に変更）
- オークション間で商品をコピー
- 商品テンプレート（よく使う商品をテンプレート保存）
- CSVから商品インポート
- 開始価格に基づく商品順序の自動提案
- リアルタイムコラボレーション（複数の主催者が同時編集）

---

## まとめ

このプランは、オークション管理ワークフローを「複雑な商品必須作成プロセス」から「柔軟な段階的開示モデル」へ変革します。

**主なメリット**:
- オークション作成が50%高速化（クリック数削減）
- 明確なワークフロー（作成 → 編集 → 商品追加）
- 統一された商品管理インターフェース
- ドラッグ&ドロップの視覚フィードバック向上
- モバイル体験の改善
- 既存データとAPIへの破壊的変更なし

**リスク軽減**:
- 段階的ロールアウトで徐々に採用
- 後方互換性のあるAPI変更
- 各フェーズでの包括的テスト
- 問題発生時の明確なロールバック戦略
