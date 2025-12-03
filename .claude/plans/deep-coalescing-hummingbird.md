# 管理者ダッシュボード実装計画

## 概要

リアルタイムオークションシステムの管理者ダッシュボード画面（`/admin/dashboard`）を実装します。これは管理者（system_adminおよびauctioneerロール）がログイン後に最初にアクセスするページで、システム統計、最近のアクティビティ、クイックアクションボタンを一目で確認できます。

**現状**: ダッシュボードルートは存在しますが、ユーザー情報を表示するプレースホルダーのみです。

**目標**: 以下を含む完全に機能するダッシュボードを作成：
- リアルタイムシステム統計（4つのカード）
- 最近のアクティビティフィード（最新の入札、新規入札者、終了したオークション）
- ロールベースのクイックアクションボタン
- 全デバイス対応のレスポンシブデザイン

## 要件概要

### システム統計カード
1. **開催中オークション数** - 現在開催中のオークション数
2. **本日の入札数** - 本日の入札総数
3. **登録入札者数** - システム内のアクティブな入札者総数
4. **ポイント流通量** - 全入札者に配布された総ポイント数

### 最近のアクティビティ
1. **最新の入札**（直近5件） - オークション名、入札者名、価格、タイムスタンプ
2. **新規入札者**（直近5件） - メールアドレス、表示名、登録日（system_adminのみ）
3. **終了したオークション**（直近5件） - タイトル、落札者、最終価格、終了時刻

### クイックアクション（ロールベース）
1. **新規オークション作成** - 全管理者が利用可能
2. **新規入札者作成** - system_adminのみ
3. **ポイント付与** - system_adminのみ

## 実装アプローチ

### フェーズ1: バックエンドAPI実装

**作成する新規ファイル:**
- `backend/internal/domain/dashboard.go` - ダッシュボードデータ構造
- `backend/internal/repository/dashboard_repository.go` - 統計情報のデータベースクエリ
- `backend/internal/service/dashboard_service.go` - ダッシュボードデータのビジネスロジック
- `backend/internal/handler/dashboard_handler.go` - ダッシュボードエンドポイントのHTTPハンドラー

**変更する既存ファイル:**
- `backend/cmd/api/main.go` - 管理者保護グループにダッシュボードルートを追加

**APIエンドポイント:**
1. `GET /api/admin/dashboard/stats` - 4つの統計情報を全て返す
2. `GET /api/admin/dashboard/activities` - 最近のアクティビティを返す（ロールベースフィルタリング）

**データベースクエリ:**
- 開催中オークション: `SELECT COUNT(*) FROM auctions WHERE status = 'active'`
- 本日の入札: `SELECT COUNT(*) FROM bids WHERE bid_at >= CURRENT_DATE`
- 入札者総数: `SELECT COUNT(*) FROM bidders WHERE status = 'active'`
- ポイント流通量: `SELECT SUM(total_points) FROM bidder_points`
- 最新の入札: `bids`, `items`, `bidders` テーブルをJOIN、`bid_at DESC` でソート、5件取得
- 新規入札者: `bidders` から選択、`created_at DESC` でソート、5件取得
- 終了したオークション: `auctions`, `items`, `bidders` をJOIN、`status = 'ended'` でフィルタ、5件取得

**ロールベースロジック:**
- `new_bidders` データはsystem_adminロールにのみ返される
- Auctioneerには空配列が返されるか、フィールドが省略される

### フェーズ2: フロントエンド状態管理

**作成する新規ファイル:**
- `frontend/src/services/api/dashboardApi.js` - ダッシュボードエンドポイント用APIクライアント
- `frontend/src/stores/dashboardStore.js` - ダッシュボード状態管理用Piniaストア
- `frontend/src/utils/timeFormatter.js` - 相対時刻フォーマット用ユーティリティ（例: "3分前"）

**Piniaストア構造:**
```javascript
{
  state: {
    stats: { activeAuctions, todayBids, totalBidders, totalPoints },
    activities: { recentBids, newBidders, endedAuctions },
    loading: false,
    error: null
  },
  actions: {
    fetchStats(),
    fetchActivities(),
    fetchAll()
  }
}
```

### フェーズ3: フロントエンドUIコンポーネント

**作成する新規コンポーネント:**
- `frontend/src/components/admin/StatsCard.vue` - アイコン、タイトル、値を持つ再利用可能な統計カード
- `frontend/src/components/admin/RecentBidsList.vue` - 相対時刻付きの最新入札リスト
- `frontend/src/components/admin/NewBiddersList.vue` - 新規入札者リスト（system_adminのみ）
- `frontend/src/components/admin/EndedAuctionsList.vue` - 最近終了したオークションリスト
- `frontend/src/components/admin/QuickActions.vue` - ロールベースアクションボタン

**変更する既存ファイル:**
- `frontend/src/views/admin/DashboardView.vue` - プレースホルダーを完全なダッシュボードレイアウトに置き換え

**コンポーネントレイアウト:**
```
DashboardView.vue
├── ヘッダー（"ダッシュボード"）
├── 統計グリッド（4つのStatsCardコンポーネント）
├── コンテンツグリッド（デスクトップで3カラム、モバイル/タブレットで1-2カラム）
│   ├── RecentBidsList
│   ├── EndedAuctionsList
│   └── サイドバー
│       ├── QuickActions
│       └── NewBiddersList（system_adminの場合）
```

### フェーズ4: スタイリングとレスポンシブ対応

**レスポンシブブレークポイント:**
- **モバイル（< 768px）**: 1カラムレイアウト、統計カードを縦に積み重ね
- **タブレット（768px - 1279px）**: 統計は2カラムグリッド、コンテンツは混合レイアウト
- **デスクトップ（≥ 1280px）**: 統計は4カラムグリッド、コンテンツは3カラムレイアウト

**デザインシステム:**
- Shadcn VueコンポーネントとTailwind CSSを使用
- カラースキーム: プライマリアクションは青、背景はグレー
- 統計カードにホバーエフェクト（シャドウの高さ変更）
- データ取得中はローディングスケルトン画面
- リトライボタン付きのエラー状態

### フェーズ5: テストと検証

**バックエンドテスト:**
- リポジトリ層: 各統計クエリが正しいカウントを返すことをテスト
- サービス層: ロールベースデータフィルタリングをテスト
- ハンドラー層: JWT認証、ロール検証、レスポンス形式をテスト

**フロントエンドテスト:**
- コンポーネントテスト: 統計カードが正しいpropsでレンダリングされることを確認
- ストアテスト: API呼び出しと状態更新を確認
- 統合テスト: ロールベースのUIレンダリングを確認

## 重要なファイル

### 作成必須（バックエンド）:
1. `backend/internal/handler/dashboard_handler.go` - コアAPIロジック
2. `backend/internal/repository/dashboard_repository.go` - データベース集計
3. `backend/internal/service/dashboard_service.go` - ビジネスロジック
4. `backend/internal/domain/dashboard.go` - データ構造

### 作成必須（フロントエンド）:
1. `frontend/src/stores/dashboardStore.js` - 状態管理
2. `frontend/src/components/admin/StatsCard.vue` - 統計表示
3. `frontend/src/services/api/dashboardApi.js` - APIクライアント
4. `frontend/src/utils/timeFormatter.js` - 時刻フォーマットユーティリティ

### 変更必須:
1. `backend/cmd/api/main.go` - ダッシュボードルートを追加
2. `frontend/src/views/admin/DashboardView.vue` - ダッシュボードUI完成

## 実装手順

### ステップ1: バックエンド基盤（3-4時間）
1. ダッシュボードデータのドメイン構造体を作成（`dashboard.go`）
2. 各統計のリポジトリメソッドを実装（`dashboard_repository.go`）
3. ロールベースフィルタリング付きサービス層を実装（`dashboard_service.go`）
4. 2つのエンドポイント用HTTPハンドラーを作成（`dashboard_handler.go`）
5. `main.go` の管理者保護グループにルートを登録
6. curl/Postmanでテスト

### ステップ2: フロントエンド状態層（2-3時間）
1. APIクライアント関数を作成（`dashboardApi.js`）
2. state、getters、actionsを持つPiniaストアを作成（`dashboardStore.js`）
3. 時刻フォーマッターユーティリティを作成（`timeFormatter.js`）
4. API統合をテスト

### ステップ3: フロントエンドUIコンポーネント（4-5時間）
1. `StatsCard.vue` 構築 - アイコン、タイトル、大きな数値表示
2. `RecentBidsList.vue` 構築 - 相対時刻付きテーブル/リスト
3. `NewBiddersList.vue` 構築 - ロールに基づく条件付きレンダリング
4. `EndedAuctionsList.vue` 構築 - 落札者と最終価格
5. `QuickActions.vue` 構築 - ロールベースボタン表示
6. `DashboardView.vue` 更新 - 全コンポーネントを統合、グリッドレイアウト追加

### ステップ4: スタイリングと仕上げ（2-3時間）
1. レスポンシブブレークポイントを実装（モバイル/タブレット/デスクトップ）
2. データ取得中のローディングスケルトンを追加
3. リトライボタン付きエラーハンドリングUIを追加
4. ホバーエフェクトとトランジションを追加
5. アクセシビリティを確認（キーボードナビゲーション、スクリーンリーダー）

### ステップ5: テスト（2-3時間）
1. バックエンドのリポジトリ/サービス/ハンドラーのユニットテストを記述
2. フロントエンドコンポーネントテストを記述
3. ロール別の手動テスト（system_admin vs auctioneer）
4. デバイス別の手動テスト（レスポンシブ）

**総推定時間**: 13-18時間

## 成功基準

- [ ] ダッシュボード読み込み時に4つの統計が正しく表示される
- [ ] 最近のアクティビティが各最大5件表示される
- [ ] 新規入札者リストはsystem_adminのみに表示される
- [ ] クイックアクションボタンが動作し、ロール権限を尊重する
- [ ] データ取得中にローディング状態が表示される
- [ ] リトライオプション付きのエラー状態が表示される
- [ ] モバイル/タブレット/デスクトップでレスポンシブデザインが機能する
- [ ] 全てのリンクが正しいページに遷移する
- [ ] ダッシュボードに戻った時にデータが更新される
- [ ] 全てのテストが合格する

## セキュリティ考慮事項

1. **認証**: 全エンドポイントで有効なJWTトークンが必要
2. **認可**: ダッシュボードはsystem_adminとauctioneerのみアクセス可能（bidderは403）
3. **データフィルタリング**: 新規入札者データはauctioneerロールには除外
4. **PII漏洩防止**: 表示名のみ表示、機密性の高い入札者UUIDは公開しない
5. **レート制限**: 統計エンドポイントへのレート制限追加を検討（将来）

## パフォーマンス考慮事項

1. **データベースクエリ**: 全てCOUNT/SUM集計とLIMITを使用して効率化
2. **キャッシュ戦略**: 統計は30-60秒キャッシュ可能（将来の機能拡張）
3. **遅延読み込み**: コンポーネントはマウント時にデータを読み込み、一度に全て読み込まない
4. **ページネーション**: アクティビティリストは5件に制限してレスポンスを小さく保つ

## 将来の機能拡張

1. **自動更新**: 30秒ごとにWebSocketまたはポーリングでリアルタイム更新
2. **チャート/グラフ**: 入札トレンド、オークション履歴の時系列視覚化
3. **日付フィルター**: 日付範囲による統計フィルタリング
4. **エクスポート**: 統計のCSV/PDFレポート
5. **カスタマイズ**: 管理者が表示する統計を選択可能に
6. **通知**: 新しいアクティビティのバッジインジケーター

## 参考資料

- 要件: `docs/screen_list.md`（35-59行目）
- データベーススキーマ: `docs/database_definition.md`
- 既存パターン: `frontend/src/views/admin/AuctionListView.vue`
- 認証ストア: `frontend/src/stores/auth.js`（`isSystemAdmin`、`isAuctioneer`を含む）
- APIルート: `backend/cmd/api/main.go`
