# Plan: システム管理者編集画面の実装

管理者編集画面（`/admin/admins/:id/edit`）を実装します。バックエンドに詳細取得・更新APIを追加し、フロントエンドの既存プレースホルダー（`AdminEditView.vue`）を完全な編集画面に置き換えます。登録画面（`AdminRegisterView.vue`）のフォーム構成を再利用し、一貫したUIを維持します。

## Steps

1. **バックエンドAPIエンドポイント追加**: [`backend/internal/handler/admin_handler.go`](backend/internal/handler/admin_handler.go) に `GetAdmin` と `UpdateAdmin` ハンドラーを追加し、[`cmd/api/main.go`](backend/cmd/api/main.go) でルーティング（`GET/PUT /admin/admins/:id`）を登録

2. **リポジトリ・サービス層の拡張**: [`internal/repository/postgres/admin_repository.go`](backend/internal/repository/postgres/admin_repository.go) に `UpdateAdmin` メソッドを実装し、[`internal/service/admin_service.go`](backend/internal/service/admin_service.go) でバリデーション付き更新ロジックを追加（メール重複チェック、自己ロール変更禁止）

3. **ドメインモデル追加**: [`internal/domain/admin.go`](backend/internal/domain/admin.go) に `AdminUpdateRequest` 構造体を定義（Email, DisplayName, Role, Status, Password（任意）フィールド）

4. **フロントエンドAPI・Store拡張**: [`frontend/src/services/api/admin.js`](frontend/src/services/api/admin.js) に `getAdmin`, `updateAdmin` 関数を追加し、[`stores/admin.js`](frontend/src/stores/admin.js) に対応するアクションを追加

5. **編集画面コンポーネント実装**: [`AdminEditView.vue`](frontend/src/views/admin/AdminEditView.vue) を [`AdminRegisterView.vue`](frontend/src/views/admin/AdminRegisterView.vue) を参考に完全実装（フォーム表示、パスワード変更オプション、状態変更、更新処理）

## Further Considerations

1. **自己編集の制限**: 自分自身のロール変更・削除を禁止するか、警告表示で確認するか → 禁止（サーバー側で拒否）を推奨
2. **最後のsystem_admin保護**: 唯一のsystem_adminを削除/停止できないようサーバー側でチェックするか → チェックを実装すべき
3. **削除機能の実装範囲**: 論理削除（status='deleted'）をこの画面で実装するか、削除ボタンは一覧画面に配置するか → 編集画面下部に削除ボタンを配置し、確認モーダル表示を推奨
