# Plan: 入札者オークション開催中画面のUI/UX改善（最終版）

**音声通知の追加**（Web Audio APIでプログラム生成）、**参加者数表示**（初期値取得対応）、**ポイント表示のコンパクト化**、**WebSocket切断時警告**、**オークションカラーの統一**の5つの改善を実装します。

## Steps

### 1. 音声通知用Composableを新規作成
- `frontend/src/composables/useAudioNotification.js` を作成
- Web Audio API の `OscillatorNode` を使用
- 3種類のプログラム生成音を実装:
  - 価格開示時: 高めのビープ2回
  - 入札成功時: 短い上昇音
  - 他者入札時: 低めの単音

### 2. ストアに参加者数管理と音声通知を追加
- `frontend/src/stores/bidderAuctionLive.js` に変更
- `participantCount` ref を追加
- `onWebSocketConnected()` 内で `participants:list` イベントをsubscribe
- `participant:joined` / `participant:left` イベントハンドラーを実装
- `onPriceOpened()` と `onBidPlaced()` に音声再生を組み込み

### 3. ヘッダーに参加者数を表示
- `frontend/src/views/bidder/BidderAuctionLive.vue` のヘッダー部分を修正
- WebSocket接続状態の左側に人型アイコンと参加者数を追加
- 例: 👥 12人参加中

### 4. PointsDisplayをコンパクト化
- `frontend/src/components/bidder/PointsDisplay.vue` を修正
- 利用可能ポイントのみを1行のインラインバッジ形式で表示
- 削除対象:
  - 3カラムグリッド（合計/利用可能/予約済み）
  - プログレスバー
  - 警告メッセージ

### 5. WebSocket切断時オーバーレイを追加
- `frontend/src/views/bidder/BidderAuctionLive.vue` に追加
- 表示条件: `wsConnected === false && !wsReconnecting`
- 固定オーバーレイで再接続ボタンと状態を大きく表示

### 6. オークションカラー（gold）を適用
- `frontend/src/components/bidder/BidPanel.vue`:
  - 入札ボタンに `bg-auction` (#D4AF37) を適用
  - 勝者バッジに `text-auction` / `border-auction` を適用
- `frontend/src/components/bidder/BidderBidHistory.vue`:
  - 勝者入札行にオークションカラーを適用

## Further Considerations

### 1. `participants:list` イベントの実装確認
- バックエンド `backend/internal/ws/hub.go` の `GetActiveParticipants()` は既に存在
- WebSocket subscribe時に初期参加者リストを返す処理がサーバー側で実装されているか確認が必要
- 未実装の場合はバックエンド側の対応も計画に追加

### 2. 音声再生のブラウザ制限対応
- ユーザー操作（入札ボタンクリック等）なしでは音声再生がブロックされる
- 対応案:
  - 初回入札時に `AudioContext.resume()` を呼ぶ
  - または画面表示時に「音声通知を有効化」する非表示のユーザー操作を検討

## 関連ファイル

| ファイル | 変更内容 |
|----------|----------|
| `frontend/src/composables/useAudioNotification.js` | 新規作成 |
| `frontend/src/stores/bidderAuctionLive.js` | 参加者数管理・音声通知追加 |
| `frontend/src/views/bidder/BidderAuctionLive.vue` | ヘッダー改修・オーバーレイ追加 |
| `frontend/src/components/bidder/PointsDisplay.vue` | コンパクト化 |
| `frontend/src/components/bidder/BidPanel.vue` | オークションカラー適用 |
| `frontend/src/components/bidder/BidderBidHistory.vue` | オークションカラー適用 |
