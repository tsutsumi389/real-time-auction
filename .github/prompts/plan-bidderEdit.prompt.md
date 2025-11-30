## Plan: 入札者編集機能の実装

system_admin が入札者（bidder）の基本情報を編集する機能を実装。既存の入札者一覧・登録機能のパターンを踏襲し、詳細取得API・更新APIを新規追加、Vue編集画面を作成する。

### Steps

1. **ドメイン層に編集用構造体を追加** - `internal/domain/bidder.go` に `BidderUpdateRequest`（email, display_name, password）と `BidderDetailResponse`（ポイント情報含む）を追加

2. **リポジトリに詳細取得・更新メソッドを追加** - `internal/repository/bidder_repository.go` に `FindByIDWithPoints()`, `UpdateBidder()` インターフェース定義、`postgres/bidder_repository.go` に実装

3. **サービス層にビジネスロジックを追加** - `internal/service/bidder_service.go` に `GetBidderDetail()`, `UpdateBidder()` を追加（メール重複チェック、パスワード空欄時はスキップ）

4. **ハンドラー・ルーティングを追加** - `internal/handler/bidder_handler.go` に `GetBidderByID()`, `UpdateBidder()` ハンドラー追加、`cmd/api/main.go` に `GET /api/admin/bidders/:id`, `PUT /api/admin/bidders/:id` ルート追加

5. **フロントエンドAPIサービス更新** - `src/services/bidder.js` に `getBidderById()`, `updateBidder()` 関数追加

6. **入札者編集画面を作成** - `src/views/admin/BidderEditView.vue` を新規作成（`BidderRegistrationView.vue` をベースにフォーム構成）、`router/index.js` に `/admin/bidders/:id/edit` ルート追加

### Further Considerations

1. **パスワード変更の確認ダイアログ** - 現在のパスワード入力を必須にするか？ / 管理者権限で確認なしに変更可とするか（推奨）
2. **メールアドレス変更時の通知** - 変更後に旧・新メールアドレスへ通知メールを送信するか？ / Phase2以降で検討（推奨）
3. **状態変更の統合** - 既存の `UpdateBidderStatus` を編集画面に統合するか？ / 一覧画面の機能を維持し編集画面では参照のみ（推奨）
