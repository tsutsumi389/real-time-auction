# オークション開催中画面 - 入札者用（2.3.2）実装プラン

**作成日**: 2025年11月27日
**対象画面**: オークション開催中画面 - 入札者用（Bidder Auction Live）
**パス**: `/bidder/auctions/:id/live`
**権限**: bidder（入札者）
**優先度**: 最高（フェーズ1: リアルタイムオークションのコア機能）

---

## 1. 概要

### 1.1 目的
入札者がオークションに参加し、リアルタイムで価格の開示を確認しながら入札を行うための画面です。主催者が開示した価格でのみ入札可能なオークションモデルを実現し、ポイント残高を確認しながら安全に入札できる仕組みを提供します。

**重要な設計方針**:
- **価格固定型入札**: 入札者は現在開示されている価格でのみ入札可能（自由価格入札は不可）
- **ポイント残高の可視化**: 利用可能ポイント、予約済みポイントをリアルタイム表示
- **リアルタイム更新**: WebSocketによる価格開示、他者の入札状況の即時反映
- **シンプルなUX**: 「入札」ボタン1つで完結する直感的な操作
- **入札制限の明示**: ポイント不足、既に入札済み、商品未開始などの状態を明確に表示

### 1.2 対象ユーザー
- **bidder（入札者）**: オークションに参加し、開示された価格で入札を行うユーザー

### 1.3 主要機能
1. オークション・商品情報のリアルタイム表示
2. 現在価格の大きな表示（主催者が開示した最新価格）
3. ポイント残高の表示（利用可能/予約済み/合計）
4. 入札ボタン（現在価格での入札を1クリックで実行）
5. 入札履歴の表示（自分の入札と他者の入札を区別）
6. 商品画像・動画の表示（スライダー形式）
7. 参加者数のリアルタイム表示
8. 入札成功・失敗の即時フィードバック（トースト通知）
9. WebSocket接続状態の表示と自動再接続

---

## 2. 画面レイアウト

### 2.1 PC/タブレット画面構成（3カラムレイアウト）

```
┌──────────────────────────────────────────────────────────────────┐
│ [ヘッダー: オークション名 | ステータス | 接続状態 | 退出ボタン]    │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────┐  ┌────────────┐  ┌─────────────────────┐       │
│  │ 商品情報    │  │ 現在価格    │  │ ポイント残高        │       │
│  │ロット番号   │  │ ¥150,000   │  │ 利用可能: 200,000pt │       │
│  │商品名       │  │            │  │ 予約済み: 50,000pt  │       │
│  │説明文       │  │[入札ボタン] │  │ 合計: 250,000pt     │       │
│  └────────────┘  └────────────┘  └─────────────────────┘       │
│                                                                  │
│  ┌────────────┐  ┌────────────┐  ┌─────────────────────┐       │
│  │ 商品画像    │  │ 入札履歴    │  │ 参加者一覧          │       │
│  │[スライダー] │  │ 1. 自分     │  │ 参加者数: 12人      │       │
│  │            │  │    150,000  │  │ オンライン: 8人     │       │
│  │  [<] [>]   │  │ 2. 他の人   │  │                     │       │
│  │            │  │    150,000  │  │ [オンライン参加者]  │       │
│  └────────────┘  │ 3. 自分     │  └─────────────────────┘       │
│                  │    140,000  │                                │
│                  └────────────┘                                │
│                                                                  │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ 商品一覧（タブ形式）                                  │       │
│  │ [ロット1: 進行中] [ロット2: 待機中] [ロット3: 終了]  │       │
│  └──────────────────────────────────────────────────────┘       │
└──────────────────────────────────────────────────────────────────┘
```

### 2.2 モバイル画面構成（縦並びレイアウト）

```
┌──────────────────────┐
│ [ヘッダー]            │
│ オークション名        │
│ [進行中] [接続中]     │
├──────────────────────┤
│ ポイント残高         │
│ 利用可能: 200,000pt  │
│ 予約済み: 50,000pt   │
├──────────────────────┤
│ 現在価格             │
│ ¥150,000            │
│ [入札ボタン]         │
├──────────────────────┤
│ 商品情報             │
│ ロット1: 商品名      │
│ 説明文...            │
├──────────────────────┤
│ 商品画像             │
│ [スライダー]         │
├──────────────────────┤
│ 入札履歴             │
│ 1. 自分 150,000     │
│ 2. 他の人 150,000   │
├──────────────────────┤
│ 参加者一覧           │
│ 参加者数: 12人       │
└──────────────────────┘
```

### 2.3 デザイン要件

- **背景**: グラデーション（グレー系、bg-gray-50）
- **カード**: 白背景、影付き、角丸（rounded-lg）
- **現在価格**: 大きく目立つ表示（text-4xl以上、金色の強調 text-auction-gold）
- **入札ボタン**: 幅いっぱい、緑系（bg-auction-green）、ホバー時に濃くなる
- **ポイント表示**: 色分け（利用可能: 緑、予約済み: 黄色、合計: 青）
- **入札履歴**: 自分の入札は背景色を変える（bg-blue-50）、他者の入札は白背景
- **接続状態**: ドット（緑: 接続中、赤: 切断、黄: 再接続中）とテキスト
- **トースト通知**: 成功（緑）、エラー（赤）、情報（青）
- **アニメーション**: 価格更新時にアニメーション（animate-price-update）

---

## 3. 入力項目とバリデーション

この画面には自由入力フィールドはなく、入札は「入札ボタン」のクリックのみで実行されます。

### 3.1 入札ボタン

| 項目 | 内容 |
|------|------|
| **ラベル** | 入札する（現在価格で入札） |
| **入力形式** | ボタンクリック |
| **必須/任意** | 任意（条件を満たす場合のみ有効化） |
| **実行条件** | 1. 商品が開始されている<br>2. 商品が終了していない<br>3. 現在価格が開示されている<br>4. 利用可能ポイント ≥ 現在価格<br>5. 自分が最新の入札者でない |

**バリデーションルール（クライアント側）:**
1. 商品未開始の場合: ボタン無効化、「商品が開始されていません」
2. 商品終了済みの場合: ボタン無効化、「商品は終了しました」
3. ポイント不足の場合: ボタン無効化、「ポイントが不足しています」
4. 既に最新の入札者の場合: ボタン無効化、「既に入札済みです」
5. 現在価格が未開示の場合: ボタン無効化、「価格が開示されていません」

**バリデーションルール（サーバー側）:**
1. 商品が存在しない: 404エラー、「商品が見つかりません」
2. 商品が開始されていない: 400エラー、「商品が開始されていません」
3. 商品が終了済み: 400エラー、「商品は既に終了しています」
4. 利用可能ポイント不足: 400エラー、「ポイントが不足しています」
5. 価格が一致しない: 409エラー、「価格が変更されました。最新の価格をご確認ください」
6. 同時入札（競合）: 409エラー、「他の入札者が先に入札しました。再度お試しください」

**バリデーションタイミング:**
- クライアント: ボタンクリック前（ボタンの有効/無効を制御）
- サーバー: 入札リクエスト受信時（最終的な検証）

---

## 4. 処理フロー

### 4.1 画面初期化フロー

```
[画面マウント]
   ↓
[JWT認証トークン確認]
   ↓
   ├→ トークンなし → [ログイン画面へリダイレクト]
   ↓
   └→ トークンあり
      ↓
   [1. オークション詳細取得（GET /api/auctions/:id）]
      ↓
   [2. ポイント残高取得（GET /api/bidder/points）]
      ↓
   [3. 商品一覧を取得、最初のアクティブ商品を選択]
      ↓
   [4. 選択商品の入札履歴取得（GET /api/bidder/items/:id/bids）]
      ↓
   [5. WebSocket接続確立（ws://localhost/ws?token=xxx&auction_id=xxx）]
      ↓
   [6. WebSocketイベントハンドラー登録]
      - price:opened（価格開示）
      - bid:placed（入札発生）
      - item:started（商品開始）
      - item:ended（商品終了）
      - auction:ended（オークション終了）
      - auction:cancelled（オークション中止）
      ↓
   [7. 画面表示完了]
```

### 4.2 入札実行フロー

```
[入札ボタンクリック]
   ↓
[クライアント側バリデーション]
   ↓
   ├→ 検証失敗 → [エラートースト表示] → [終了]
   ↓
   └→ 検証成功
      ↓
   [楽観的UI更新（ボタン無効化、ローディング表示）]
      ↓
   [POST /api/bidder/items/:id/bid]
   リクエストボディ: { "price": 150000 }
      ↓
   ┌─────────────────────┐
   │ サーバー側処理      │
   │ 1. JWT認証          │
   │ 2. バリデーション   │
   │ 3. Redis分散ロック  │
   │ 4. トランザクション │
   │    - INSERT bid     │
   │    - UPDATE points  │
   │    - 他の入札を     │
   │      is_winning=false│
   │ 5. Redis Pub/Sub発行│
   │ 6. レスポンス返却   │
   └─────────────────────┘
      ↓
   ├→ 成功（200 OK）
   │  ↓
   │  [成功トースト表示: "入札しました"]
   │  ↓
   │  [WebSocketで全参加者に通知]
   │  イベント: bid:placed
   │  ↓
   │  [入札履歴を先頭に追加]
   │  [ポイント残高を更新（available減、reserved増）]
   │  [入札ボタン無効化（自分が最新入札者）]
   │  ↓
   │  [終了]
   │
   └→ 失敗（400/409/500）
      ↓
      [エラートースト表示]
      - 400: バリデーションエラー
      - 409: 競合エラー（価格変更、同時入札）
      - 500: サーバーエラー
      ↓
      [UI状態を元に戻す（ボタン有効化）]
      ↓
      [終了]
```

### 4.3 価格開示受信フロー（WebSocket）

```
[WebSocketイベント受信: price:opened]
ペイロード: {
  "item_id": "uuid",
  "price": 160000,
  "previous_price": 150000,
  "disclosed_at": "2025-11-27T10:30:00Z"
}
   ↓
[現在商品の価格を更新]
   ↓
[価格変更アニメーション表示（animate-price-update）]
   ↓
[入札ボタン状態を再評価]
   ↓
   ├→ 自分が前回の入札者の場合
   │  ↓
   │  [自分の入札をis_winning=falseに更新]
   │  [予約済みポイントを利用可能に戻す（UI表示のみ）]
   │  [入札ボタン有効化]
   │
   └→ 自分が前回の入札者でない場合
      ↓
      [入札ボタン状態を再評価]
      - ポイント不足の場合は無効化
      - ポイント足りる場合は有効化
   ↓
[情報トースト表示: "価格が更新されました: ¥160,000"]
   ↓
[終了]
```

### 4.4 他者の入札受信フロー（WebSocket）

```
[WebSocketイベント受信: bid:placed]
ペイロード: {
  "bid": {
    "id": 123,
    "item_id": "uuid",
    "bidder_id": "uuid",
    "bidder_name": "入札者B",
    "price": 150000,
    "is_winning": true,
    "bid_at": "2025-11-27T10:30:00Z"
  }
}
   ↓
[入札が自分のものかチェック]
   ↓
   ├→ 自分の入札の場合
   │  ↓
   │  [既にクライアント側で処理済みなのでスキップ]
   │  ↓
   │  [終了]
   │
   └→ 他者の入札の場合
      ↓
      [入札履歴に追加（先頭に挿入）]
      ↓
      [自分の前回入札があればis_winning=falseに更新]
      ↓
      [自分の予約済みポイントを利用可能に戻す（UI表示のみ）]
      ↓
      [入札ボタン有効化（再入札可能に）]
      ↓
      [情報トースト表示: "他の入札者が入札しました"]
      ↓
      [終了]
```

### 4.5 商品終了受信フロー（WebSocket）

```
[WebSocketイベント受信: item:ended]
ペイロード: {
  "item": {
    "id": "uuid",
    "winner_id": "uuid",
    "winner_name": "入札者A",
    "final_price": 150000,
    "ended_at": "2025-11-27T10:35:00Z"
  }
}
   ↓
[商品のステータスを'ended'に更新]
   ↓
[入札ボタン無効化]
   ↓
[落札者が自分かチェック]
   ↓
   ├→ 自分が落札者の場合
   │  ↓
   │  [成功トースト表示: "落札しました！"]
   │  [予約済みポイントを消費（reserved減）]
   │  [落札商品として表示（金色のバッジ）]
   │
   └→ 他者が落札者の場合
      ↓
      [情報トースト表示: "商品が終了しました。落札者: ○○さん"]
      [予約済みポイントを利用可能に戻す（reserved → available）]
   ↓
[次の待機中商品を自動選択]
   ↓
   ├→ 待機中商品がある場合
   │  ↓
   │  [次の商品を選択]
   │  [次の商品の入札履歴を取得]
   │
   └→ 待機中商品がない場合
      ↓
      [商品一覧を表示（すべて終了）]
   ↓
[終了]
```

---

## 5. セキュリティ要件

### 5.1 認証・認可

1. **JWT認証**: すべてのAPIリクエストに`Authorization: Bearer <token>`ヘッダーを含める
2. **WebSocket認証**: 接続時にクエリパラメータ`?token=xxx`でJWTトークンを渡す
3. **入札者ロール検証**: サーバー側で`user_type=bidder`を検証
4. **トークン有効期限**: 24時間、期限切れの場合はログイン画面へリダイレクト
5. **CORS設定**: 開発環境では`*`、本番環境では許可されたオリジンのみ

### 5.2 入札の整合性

1. **Redis分散ロック**: 同時入札を防ぐため、入札処理中はRedisでロック（キー: `bid:lock:item:{item_id}`）
2. **トランザクション**: 入札レコード作成とポイント更新を同一トランザクション内で実行
3. **楽観的ロック**: `current_price`を条件にUPDATE、変更されていたら409エラー
4. **べき等性**: 同一bidder_idと同一priceの入札は重複排除（既存の場合は200 OK）

### 5.3 ポイント管理の安全性

1. **CHECKコンストラクション**: PostgreSQLレベルで`CHECK (available_points + reserved_points <= total_points)`
2. **非同期整合性**: ポイント操作後にRedisにキャッシュ（キー: `bidder:points:{bidder_id}`）
3. **監査ログ**: すべてのポイント操作を`point_history`テーブルに記録（before/after値）
4. **負数防止**: `available_points`が負にならないようサーバー側でチェック

### 5.4 XSS/CSRF対策

1. **XSS対策**: Vueの自動エスケープ（`{{ }}`）により、ユーザー入力は自動的にエスケープ
2. **CSRF対策**: JWTトークン認証のため、CSRFトークンは不要（ステートレス）
3. **Content Security Policy**: HTTPヘッダーで`script-src`を制限（本番環境）

### 5.5 WebSocket接続の保護

1. **認証済み接続のみ**: 認証されていない接続はWebSocketハンドシェイク時に拒否
2. **Ping/Pong**: 30秒間隔でPingを送信、60秒間Pongがない場合は切断
3. **再接続制限**: 最大5回まで再接続試行、失敗した場合はエラー表示

---

## 6. データベース設計

### 6.1 既存テーブル（変更なし）

既存のデータベーススキーマを使用します。以下は関連するテーブルの概要です。

**auctions テーブル（オークション）**

| カラム名 | 型 | 制約 | 説明 |
|---------|-----|------|------|
| id | UUID | PRIMARY KEY | オークションID |
| title | VARCHAR(255) | NOT NULL | オークション名 |
| description | TEXT |  | 説明 |
| status | VARCHAR(20) | NOT NULL, CHECK | ステータス（active/ended/cancelled） |
| started_at | TIMESTAMPTZ |  | 開始日時 |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 |
| updated_at | TIMESTAMPTZ | NOT NULL | 更新日時 |

**items テーブル（商品）**

| カラム名 | 型 | 制約 | 説明 |
|---------|-----|------|------|
| id | UUID | PRIMARY KEY | 商品ID |
| auction_id | UUID | NOT NULL, FK | オークションID |
| name | VARCHAR(255) | NOT NULL | 商品名 |
| description | TEXT |  | 説明 |
| lot_number | INT | NOT NULL | ロット番号 |
| starting_price | BIGINT |  | 開始価格 |
| current_price | BIGINT |  | 現在価格 |
| winner_id | UUID | FK | 落札者ID |
| started_at | TIMESTAMPTZ |  | 商品開始日時 |
| ended_at | TIMESTAMPTZ |  | 商品終了日時 |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 |
| updated_at | TIMESTAMPTZ | NOT NULL | 更新日時 |

**bids テーブル（入札）**

| カラム名 | 型 | 制約 | 説明 |
|---------|-----|------|------|
| id | BIGSERIAL | PRIMARY KEY | 入札ID |
| item_id | UUID | NOT NULL, FK | 商品ID |
| bidder_id | UUID | NOT NULL, FK | 入札者ID |
| price | BIGINT | NOT NULL | 入札価格 |
| is_winning | BOOLEAN | NOT NULL, DEFAULT false | 最新入札フラグ |
| bid_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | 入札日時 |

**bidder_points テーブル（ポイント残高）**

| カラム名 | 型 | 制約 | 説明 |
|---------|-----|------|------|
| bidder_id | UUID | PRIMARY KEY, FK | 入札者ID |
| total_points | BIGINT | NOT NULL, DEFAULT 0 | 合計ポイント |
| available_points | BIGINT | NOT NULL, DEFAULT 0 | 利用可能ポイント |
| reserved_points | BIGINT | NOT NULL, DEFAULT 0 | 予約済みポイント |
| updated_at | TIMESTAMPTZ | NOT NULL | 更新日時 |

**CHECK制約**: `CHECK (available_points + reserved_points <= total_points)`

**point_history テーブル（ポイント履歴）**

| カラム名 | 型 | 制約 | 説明 |
|---------|-----|------|------|
| id | BIGSERIAL | PRIMARY KEY | 履歴ID |
| bidder_id | UUID | NOT NULL, FK | 入札者ID |
| operation_type | VARCHAR(20) | NOT NULL | 操作種別（reserve/release/consume/refund） |
| amount | BIGINT | NOT NULL | 変動量 |
| before_available | BIGINT | NOT NULL | 変更前の利用可能ポイント |
| after_available | BIGINT | NOT NULL | 変更後の利用可能ポイント |
| before_reserved | BIGINT | NOT NULL | 変更前の予約済みポイント |
| after_reserved | BIGINT | NOT NULL | 変更後の予約済みポイント |
| reference_type | VARCHAR(20) |  | 参照タイプ（bid/refund） |
| reference_id | BIGINT |  | 参照ID（bid.id等） |
| created_at | TIMESTAMPTZ | NOT NULL | 作成日時 |

### 6.2 テストデータ仕様

**オークション1件:**
- title: "競走馬オークション2025"
- status: "active"
- 商品3件を含む

**商品3件:**
- ロット1: 開始済み、current_price=150000
- ロット2: 待機中、starting_price=200000
- ロット3: 終了済み、winner_id=bidder_aのUUID

**入札者2名:**
- bidder_a: 250000 total_points (200000 available, 50000 reserved)
- bidder_b: 300000 total_points (300000 available, 0 reserved)

**入札履歴（ロット1）:**
- bidder_a: 150000, is_winning=true
- bidder_b: 150000, is_winning=false
- bidder_a: 140000, is_winning=false

---

## 7. API仕様

### 7.1 オークション詳細取得

**エンドポイント**: `GET /api/auctions/:id`
**認証**: 不要（公開エンドポイント）
**説明**: オークションの詳細と商品一覧を取得

**リクエスト例:**
```
GET /api/auctions/550e8400-e29b-41d4-a716-446655440000
```

**レスポンス例（200 OK）:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "競走馬オークション2025",
  "description": "優秀な競走馬のオークション",
  "status": "active",
  "started_at": "2025-11-27T10:00:00Z",
  "created_at": "2025-11-20T09:00:00Z",
  "updated_at": "2025-11-27T10:00:00Z",
  "items": [
    {
      "id": "item-uuid-1",
      "name": "サンデーサイレンス産駒",
      "description": "3歳牡馬",
      "lot_number": 1,
      "starting_price": 100000,
      "current_price": 150000,
      "winner_id": null,
      "started_at": "2025-11-27T10:05:00Z",
      "ended_at": null,
      "media": [
        {
          "id": "media-uuid-1",
          "url": "https://example.com/image1.jpg",
          "type": "image",
          "display_order": 1
        }
      ]
    }
  ]
}
```

**エラーレスポンス:**
- 404: `{"error": "Auction not found"}`
- 500: `{"error": "Internal server error"}`

---

### 7.2 ポイント残高取得

**エンドポイント**: `GET /api/bidder/points`
**認証**: 必須（`Authorization: Bearer <token>`）
**説明**: 現在ログイン中の入札者のポイント残高を取得

**リクエスト例:**
```
GET /api/bidder/points
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**レスポンス例（200 OK）:**
```json
{
  "bidder_id": "bidder-uuid-1",
  "total_points": 250000,
  "available_points": 200000,
  "reserved_points": 50000,
  "updated_at": "2025-11-27T10:30:00Z"
}
```

**エラーレスポンス:**
- 401: `{"error": "Unauthorized"}` - トークンなし、または無効
- 403: `{"error": "Forbidden"}` - 入札者ロールでない
- 404: `{"error": "Points not found"}` - ポイント情報が存在しない
- 500: `{"error": "Internal server error"}`

---

### 7.3 入札実行

**エンドポイント**: `POST /api/bidder/items/:id/bid`
**認証**: 必須（`Authorization: Bearer <token>`）
**説明**: 指定した商品に現在価格で入札

**リクエスト例:**
```
POST /api/bidder/items/item-uuid-1/bid
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "price": 150000
}
```

**レスポンス例（200 OK）:**
```json
{
  "bid": {
    "id": 123,
    "item_id": "item-uuid-1",
    "bidder_id": "bidder-uuid-1",
    "price": 150000,
    "is_winning": true,
    "bid_at": "2025-11-27T10:30:00Z"
  },
  "points": {
    "total_points": 250000,
    "available_points": 100000,
    "reserved_points": 150000,
    "updated_at": "2025-11-27T10:30:00Z"
  }
}
```

**エラーレスポンス:**
- 400: `{"error": "Item not started"}` - 商品が開始されていない
- 400: `{"error": "Item already ended"}` - 商品が終了済み
- 400: `{"error": "Insufficient points"}` - ポイント不足
- 404: `{"error": "Item not found"}` - 商品が存在しない
- 409: `{"error": "Price has changed. Please check the latest price"}` - 価格が変更された
- 409: `{"error": "Another bidder placed a bid first. Please try again"}` - 同時入札で競合
- 500: `{"error": "Internal server error"}`

---

### 7.4 入札履歴取得

**エンドポイント**: `GET /api/bidder/items/:id/bids`
**認証**: 必須（`Authorization: Bearer <token>`）
**説明**: 指定した商品の入札履歴を取得（自分の入札と他者の入札を区別）

**クエリパラメータ:**
- `limit` (optional): 取得件数（デフォルト: 50、最大: 200）
- `offset` (optional): オフセット（デフォルト: 0）

**リクエスト例:**
```
GET /api/bidder/items/item-uuid-1/bids?limit=20&offset=0
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**レスポンス例（200 OK）:**
```json
{
  "total": 3,
  "bids": [
    {
      "id": 123,
      "item_id": "item-uuid-1",
      "bidder_id": "bidder-uuid-1",
      "bidder_name": "自分",
      "price": 150000,
      "is_winning": true,
      "is_own_bid": true,
      "bid_at": "2025-11-27T10:30:00Z"
    },
    {
      "id": 122,
      "bidder_id": "bidder-uuid-2",
      "bidder_name": "入札者B",
      "price": 150000,
      "is_winning": false,
      "is_own_bid": false,
      "bid_at": "2025-11-27T10:29:00Z"
    }
  ]
}
```

**エラーレスポンス:**
- 401: `{"error": "Unauthorized"}`
- 404: `{"error": "Item not found"}`
- 500: `{"error": "Internal server error"}`

---

## 8. バックエンド実装要件

### 8.1 ディレクトリ構成

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # 入札者用ルートを追加
├── internal/
│   ├── domain/
│   │   ├── bid.go               # 既存（変更なし）
│   │   └── bidder_points.go     # 新規作成
│   ├── repository/
│   │   ├── bid_repository.go    # 新規作成
│   │   └── point_repository.go  # 新規作成
│   ├── service/
│   │   ├── bid_service.go       # 新規作成（入札ビジネスロジック）
│   │   └── point_service.go     # 新規作成（ポイント管理）
│   ├── handler/
│   │   └── bid_handler.go       # 新規作成（入札APIハンドラー）
│   └── middleware/
│       └── auth.go              # 既存のRequireBidder()を使用
```

### 8.2 Domain層（ドメインモデル）

**bidder_points.go** - ポイント残高モデル

何を定義するか:
- BidderPoints構造体（bidder_id, total_points, available_points, reserved_points, updated_at）
- PointHistory構造体（ポイント履歴のレコード）
- GetPointsResponse（APIレスポンス用）

技術要件:
- GORMタグでフィールドをマッピング
- JSONタグでAPIレスポンス形式を定義

---

### 8.3 Repository層（データアクセス）

**bid_repository.go** - 入札データアクセス

何を提供するか:
- CreateBid(bid *Bid) - 入札レコード作成
- FindBidsByItemID(itemID string, limit, offset int) - 入札履歴取得
- CountBidsByItemID(itemID string) - 入札件数取得
- UpdateBidWinningStatus(itemID string, bidID int64) - is_winningフラグ更新（トリガーが既に処理）
- FindWinningBidByItemID(itemID string) - 最新の入札を取得

技術要件:
- GORMを使用してPostgreSQLにアクセス
- トランザクション対応（*gorm.DBを受け取る）

**point_repository.go** - ポイントデータアクセス

何を提供するか:
- FindPointsByBidderID(bidderID string) - ポイント残高取得
- UpdatePoints(bidderID string, availableDelta, reservedDelta int64, tx *gorm.DB) - ポイント更新
- CreatePointHistory(history *PointHistory, tx *gorm.DB) - 履歴作成

技術要件:
- CHECK制約により整合性を保証
- トランザクション内で実行

---

### 8.4 Service層（ビジネスロジック）

**bid_service.go** - 入札ビジネスロジック

何を提供するか:
- PlaceBid(itemID string, bidderID string, price int64) - 入札実行
  - バリデーション: 商品存在、開始済み、未終了、価格一致
  - Redis分散ロック取得（キー: `bid:lock:item:{itemID}`、TTL: 5秒）
  - トランザクション開始
    - ポイント確認（available_points >= price）
    - 入札レコード作成
    - ポイント更新（available - price, reserved + price）
    - ポイント履歴作成（operation_type="reserve"）
  - トランザクションコミット
  - Redisロック解放
  - Redis Pub/Sub発行（チャネル: `auction:bid`、ペイロード: 入札情報）
  - レスポンス返却

- GetBidHistory(itemID string, bidderID string, limit, offset int) - 入札履歴取得
  - 入札履歴を取得
  - `is_own_bid`フラグを追加（bidder_id == 自分のID）

技術要件:
- go-redisでRedis分散ロック（`SET key value NX EX 5`）
- GORMトランザクション（`db.Transaction()`）
- エラーハンドリング（ロック取得失敗、ポイント不足、価格変更）

**point_service.go** - ポイント管理ビジネスロジック

何を提供するか:
- GetPoints(bidderID string) - ポイント残高取得
  - Redisキャッシュから取得（キー: `bidder:points:{bidderID}`、TTL: 60秒）
  - キャッシュミスの場合はDBから取得してキャッシュ

技術要件:
- go-redisでキャッシュ管理
- JSONシリアライズ/デシリアライズ

---

### 8.5 Handler層（HTTPハンドラー）

**bid_handler.go** - 入札APIハンドラー

何を提供するか:
- PlaceBid(c *gin.Context) - POST /api/bidder/items/:id/bid
  - JWTからbidder_idを取得（`c.Get("bidder_id")`）
  - リクエストボディから価格を取得
  - サービス層のPlaceBid()を呼び出し
  - エラーハンドリング（400/404/409/500）
  - 成功レスポンス（200 OK）

- GetBidHistory(c *gin.Context) - GET /api/bidder/items/:id/bids
  - クエリパラメータ（limit, offset）を取得
  - JWTからbidder_idを取得
  - サービス層のGetBidHistory()を呼び出し
  - レスポンス返却

- GetPoints(c *gin.Context) - GET /api/bidder/points
  - JWTからbidder_idを取得
  - サービス層のGetPoints()を呼び出し
  - レスポンス返却

技術要件:
- Ginフレームワーク
- エラーレスポンスの統一（ErrorResponse構造体）
- HTTPステータスコードの適切な使用

---

### 8.6 ルート定義（main.go）

**cmd/api/main.go** - 入札者用ルート追加

追加するルート:
```go
bidderGroup := router.Group("/api/bidder")
bidderGroup.Use(middleware.AuthMiddleware(), middleware.RequireBidder())
{
    bidderGroup.GET("/points", bidHandler.GetPoints)
    bidderGroup.POST("/items/:id/bid", bidHandler.PlaceBid)
    bidderGroup.GET("/items/:id/bids", bidHandler.GetBidHistory)
}
```

技術要件:
- AuthMiddleware: JWTトークン検証
- RequireBidder: `user_type=bidder`検証
- ハンドラーのDI（依存性注入）

---

### 8.7 WebSocketイベント発行

**service/bid_service.go内** - 入札成功時にイベント発行

何をするか:
- 入札トランザクション成功後、Redis Pub/Subでイベント発行
- チャネル: `auction:bid`
- ペイロード（JSON）:
```json
{
  "type": "bid:placed",
  "auction_id": "auction-uuid",
  "item_id": "item-uuid",
  "bid": {
    "id": 123,
    "bidder_id": "bidder-uuid",
    "bidder_name": "入札者A",
    "price": 150000,
    "is_winning": true,
    "bid_at": "2025-11-27T10:30:00Z"
  }
}
```

技術要件:
- go-redisの`Publish()`メソッド
- JSONシリアライズ（`json.Marshal()`）

---

## 9. フロントエンド実装要件

### 9.1 ディレクトリ構成

```
frontend/src/
├── views/
│   └── bidder/
│       └── BidderAuctionLive.vue   # 新規作成
├── components/
│   └── bidder/
│       ├── BidPanel.vue             # 新規作成（入札パネル）
│       ├── PointsDisplay.vue        # 新規作成（ポイント表示）
│       ├── BidderBidHistory.vue     # 新規作成（入札履歴）
│       └── ItemTabs.vue             # 新規作成（商品タブ）
├── stores/
│   └── bidderAuctionLive.js         # 新規作成（Piniaストア）
├── services/
│   └── bidderBidApi.js              # 新規作成（入札API）
└── router/
    └── index.js                     # ルート追加
```

### 9.2 メインビュー（BidderAuctionLive.vue）

何をするコンポーネントか:
- オークションライブ画面の全体レイアウトを管理
- 子コンポーネント（BidPanel, PointsDisplay, BidderBidHistory等）を配置
- WebSocket接続のライフサイクル管理
- エラー表示とトースト通知

技術要件:
- Vue 3 Composition API（`<script setup>`）
- onMounted時に初期化、onUnmounted時にクリーンアップ
- レスポンシブレイアウト（lg:grid-cols-12で3カラム、モバイルは縦並び）
- Shadcn Vueコンポーネント（Card, Button, Alert, LoadingSpinner）

---

### 9.3 入札パネル（BidPanel.vue）

何をするコンポーネントか:
- 現在価格の大きな表示
- 入札ボタン（条件に応じて有効/無効化）
- 入札ボタンクリック時に親コンポーネントへイベント発火
- ボタンの状態表示（ローディング、無効化理由）

Props:
- `item` (Object): 現在の商品情報
- `points` (Object): ポイント残高
- `isLoading` (Boolean): 入札処理中フラグ

Emits:
- `bid`: 入札ボタンクリック時

技術要件:
- 入札ボタンの有効化ロジックをcomputed propertyで実装
- 無効化理由をテキストで表示（例: "ポイントが不足しています"）
- 金額フォーマット: `new Intl.NumberFormat('ja-JP').format(price)`
- アニメーション: 価格更新時に`animate-price-update`クラス

---

### 9.4 ポイント表示（PointsDisplay.vue）

何をするコンポーネントか:
- 利用可能ポイント、予約済みポイント、合計ポイントを表示
- 色分け（利用可能: 緑、予約済み: 黄、合計: 青）
- プログレスバーで視覚的に表示

Props:
- `points` (Object): ポイント残高（total, available, reserved）

技術要件:
- プログレスバー: available_points / total_points * 100%
- Tailwind CSSのutilityクラス（text-auction-green, text-yellow-600, text-auction-blue）

---

### 9.5 入札履歴（BidderBidHistory.vue）

何をするコンポーネントか:
- 入札履歴をリスト表示
- 自分の入札と他者の入札を区別（背景色）
- スクロール可能なリスト（最大高さ制限）

Props:
- `bids` (Array): 入札履歴配列
- `currentBidderId` (String): 現在のユーザーのbidder_id

技術要件:
- 自分の入札: `bg-blue-50`、他者の入札: `bg-white`
- リスト項目: 入札者名、価格、入札日時
- 日時フォーマット: `new Intl.DateTimeFormat('ja-JP', { hour: '2-digit', minute: '2-digit', second: '2-digit' }).format(new Date(bid_at))`

---

### 9.6 Piniaストア（bidderAuctionLive.js）

何を管理するか:
- State: auction, items, currentItem, bids, points, loading, error, wsConnected, wsReconnecting
- Computed: currentPrice, hasEnoughPoints, isOwnBidWinning, canBid
- Actions:
  - `initialize(auctionId)` - オークション詳細とポイント取得
  - `fetchBidHistory(itemId)` - 入札履歴取得
  - `placeBid(itemId, price)` - 入札実行
  - `selectItem(itemId)` - 商品選択
  - `connectWebSocket(token, auctionId)` - WebSocket接続
  - `disconnectWebSocket()` - WebSocket切断
  - `onPriceOpened(payload)` - 価格開示イベントハンドラー
  - `onBidPlaced(payload)` - 入札イベントハンドラー
  - `onItemEnded(payload)` - 商品終了イベントハンドラー

技術要件:
- Pinia defineStore
- WebSocketサービスのイベント登録/解除
- 楽観的UI更新（入札時にローカル状態を先に更新）
- エラーハンドリングとロールバック

---

### 9.7 API サービス（bidderBidApi.js）

何を提供するか:
- `getAuctionDetail(auctionId)` - GET /api/auctions/:id
- `getPoints()` - GET /api/bidder/points
- `placeBid(itemId, price)` - POST /api/bidder/items/:id/bid
- `getBidHistory(itemId, limit, offset)` - GET /api/bidder/items/:id/bids

技術要件:
- Axiosインスタンス（`/src/services/apiClient.js`）を使用
- 自動的にAuthorizationヘッダーを付与（インターセプター）
- エラーレスポンスをthrowしてストアで処理

---

### 9.8 ルート定義（router/index.js）

追加するルート:
```javascript
{
  path: '/bidder/auctions/:id/live',
  name: 'bidder-auction-live',
  component: () => import('@/views/bidder/BidderAuctionLive.vue'),
  meta: {
    requiresAuth: true,
    requiresBidder: true
  }
}
```

技術要件:
- ナビゲーションガード（`requiresAuth`, `requiresBidder`）
- トークン有効性チェック（`isBidderTokenValid()`）
- 無効な場合は`/bidder/login`へリダイレクト

---

## 10. 画面遷移

### 10.1 遷移フロー図

```
[入札者ログイン画面]
   ↓（ログイン成功）
[オークション一覧画面]
   ↓（「参加する」ボタンクリック）
[オークション開催中画面] ← この画面
   ↓（オークション終了）
[オークション一覧画面]
```

### 10.2 遷移ルール

| 遷移元 | トリガー | 遷移先 | 条件 |
|-------|---------|-------|------|
| オークション一覧 | 「参加する」ボタン | オークション開催中画面 | status='active' |
| オークション開催中画面 | 「退出」ボタン | オークション一覧 | なし |
| オークション開催中画面 | オークション終了イベント | オークション一覧 | auction:ended受信 |
| オークション開催中画面 | オークション中止イベント | オークション一覧 | auction:cancelled受信 |
| オークション開催中画面 | 初期化失敗 | オークション一覧 | 404エラー |
| オークション開催中画面 | トークン期限切れ | 入札者ログイン画面 | 401エラー |

### 10.3 ナビゲーションガード

**beforeEnter（ルートレベル）:**
```javascript
beforeEnter: (to, from, next) => {
  if (!isBidderTokenValid()) {
    next('/bidder/login')
  } else {
    next()
  }
}
```

**onMounted（コンポーネントレベル）:**
- オークション詳細取得
- 404エラーの場合はオークション一覧へリダイレクト
- 401エラーの場合はログイン画面へリダイレクト

---

## 11. UIの動作仕様

### 11.1 入札ボタンの状態遷移

| 状態 | 表示 | クリック可否 | 条件 |
|------|------|-------------|------|
| 通常（有効） | 「入札する」（緑色） | 可能 | すべての条件を満たす |
| ローディング | 「入札中...」（スピナー表示） | 不可 | 入札処理中 |
| ポイント不足 | 「ポイント不足」（灰色） | 不可 | available_points < current_price |
| 商品未開始 | 「商品未開始」（灰色） | 不可 | started_at == null |
| 商品終了 | 「終了済み」（灰色） | 不可 | ended_at != null |
| 既に入札済み | 「入札済み」（青色） | 不可 | 自分が最新の入札者 |
| 価格未開示 | 「価格未開示」（灰色） | 不可 | current_price == null |

### 11.2 入札成功時のフィードバック

1. **楽観的UI更新**: ボタンをローディング状態に変更
2. **APIリクエスト送信**: POST /api/bidder/items/:id/bid
3. **成功時**:
   - トースト表示（緑色）: "入札しました"
   - 入札履歴に自分の入札を追加（先頭）
   - ポイント残高を更新（available減、reserved増）
   - 入札ボタンを「入札済み」状態に変更
4. **失敗時**:
   - トースト表示（赤色）: エラーメッセージ
   - UI状態をロールバック（ボタンを元の状態に戻す）

### 11.3 価格開示時のフィードバック

1. **WebSocketイベント受信**: `price:opened`
2. **現在価格を更新**: アニメーション付き（`animate-price-update`）
3. **トースト表示**（青色）: "価格が更新されました: ¥160,000"
4. **入札ボタン状態を再評価**:
   - 自分が前回の入札者の場合: 「入札する」に戻る
   - ポイント不足の場合: 「ポイント不足」に変更
5. **予約済みポイントを利用可能に戻す**（自分の入札がis_winning=falseになった場合）

### 11.4 他者の入札受信時のフィードバック

1. **WebSocketイベント受信**: `bid:placed`（自分以外の入札）
2. **入札履歴に追加**: 先頭に挿入
3. **トースト表示**（青色）: "他の入札者が入札しました"
4. **自分の前回入札を更新**: is_winning=false、予約済みポイントを利用可能に戻す
5. **入札ボタンを有効化**: 再入札可能に

### 11.5 商品終了時のフィードバック

1. **WebSocketイベント受信**: `item:ended`
2. **商品ステータスを'ended'に更新**
3. **入札ボタン無効化**
4. **落札者チェック**:
   - 自分が落札者: トースト表示（金色）: "落札しました！"、予約済みポイントを消費
   - 他者が落札者: トースト表示（青色）: "商品が終了しました。落札者: ○○さん"、予約済みポイントを利用可能に戻す
5. **次の待機中商品を自動選択**（あれば）

### 11.6 WebSocket接続状態の表示

| 状態 | ドット色 | テキスト | 説明 |
|------|---------|---------|------|
| 接続中 | 緑 | 接続中 | WebSocket接続が確立されている |
| 切断 | 赤 | 切断 | WebSocket接続が切断された |
| 再接続中 | 黄 | 再接続中... | 自動再接続を試行中 |

**再接続ロジック:**
- 切断時: 3秒後に再接続試行
- 最大5回まで試行
- 5回失敗した場合: エラートースト表示「接続に失敗しました。ページを再読み込みしてください」

---

## 12. レスポンシブデザイン

### 12.1 ブレークポイント

| デバイス | ブレークポイント | 画面幅 | レイアウト |
|---------|----------------|-------|-----------|
| モバイル | デフォルト | ~640px | 縦並び（1カラム） |
| タブレット | sm: | 640px~ | 縦並び（1カラム） |
| PC | lg: | 1024px~ | 3カラムレイアウト |

### 12.2 モバイル対応

**レイアウト調整:**
- ヘッダー: 2行表示（1行目: タイトル、2行目: ステータスと接続状態）
- ポイント表示: 横並び3列（利用可能、予約済み、合計）
- 現在価格: text-3xl（PCはtext-5xl）
- 入札ボタン: 幅100%
- 商品画像: スライダー（1枚ずつ表示）
- 入札履歴: 最大高さ300px、スクロール
- 参加者一覧: 折りたたみ可能

**タッチ操作:**
- ボタン: 最小タップエリア44px x 44px
- スライダー: スワイプ操作対応
- リスト: スクロール可能

### 12.3 PC対応

**レイアウト調整:**
- 3カラムグリッド（4:4:4）
- 左カラム: 商品情報と画像
- 中央カラム: 入札履歴と価格履歴
- 右カラム: 入札パネルと参加者一覧
- ヘッダー: 1行表示

---

## 13. アクセシビリティ

### 13.1 キーボード操作

| 操作 | キー | 動作 |
|------|-----|------|
| 入札ボタン | Enter / Space | 入札実行 |
| 商品タブ移動 | Tab | 次の商品タブへフォーカス |
| 商品選択 | Enter | 商品を選択 |
| モーダル閉じる | Esc | モーダルを閉じる |

### 13.2 スクリーンリーダー対応

**ARIAラベル:**
- 入札ボタン: `aria-label="現在価格 150,000ポイントで入札"`
- ポイント表示: `aria-label="利用可能ポイント 200,000、予約済み 50,000、合計 250,000"`
- 接続状態: `aria-label="WebSocket 接続中"`
- 商品タブ: `aria-label="ロット1 サンデーサイレンス産駒 進行中"`

**ライブリージョン:**
- トースト通知: `role="alert" aria-live="polite"`
- 価格更新: `aria-live="polite"`（スクリーンリーダーに読み上げ）

### 13.3 フォーカス管理

**フォーカストラップ:**
- モーダル表示時: モーダル内にフォーカスを閉じ込める
- モーダル閉じる時: 前のフォーカス位置に戻す

**フォーカス表示:**
- すべてのインタラクティブ要素にfocusリング表示（`focus:ring-2 focus:ring-blue-500`）

---

## 14. テスト要件

### 14.1 バックエンドテストシナリオ

**ユニットテスト（Service層）:**
1. 入札実行: ポイント不足の場合はエラー
2. 入札実行: 価格が変更された場合はエラー
3. 入札実行: 正常系（トランザクション成功、ポイント更新、履歴作成）
4. ポイント取得: キャッシュヒット、キャッシュミス

**統合テスト（Handler層）:**
1. POST /api/bidder/items/:id/bid: 401エラー（認証なし）
2. POST /api/bidder/items/:id/bid: 403エラー（入札者ロールでない）
3. POST /api/bidder/items/:id/bid: 400エラー（ポイント不足）
4. POST /api/bidder/items/:id/bid: 409エラー（価格変更）
5. POST /api/bidder/items/:id/bid: 200 OK（正常系）
6. GET /api/bidder/points: 200 OK（正常系）
7. GET /api/bidder/items/:id/bids: 200 OK（自分の入札フラグ確認）

**並行処理テスト:**
1. 同時入札: 2人が同時に入札した場合、1人のみ成功
2. Redis分散ロック: ロック取得失敗時に409エラー

### 14.2 フロントエンドテストシナリオ

**コンポーネントテスト（Vitest）:**
1. BidPanel: ポイント不足の場合ボタン無効化
2. BidPanel: 入札ボタンクリックでemit発火
3. PointsDisplay: ポイント表示の正確性
4. BidderBidHistory: 自分の入札と他者の入札の区別

**E2Eテスト（Playwright / Cypress）:**
1. ログイン → オークション一覧 → オークション開催中画面への遷移
2. 入札ボタンクリック → 入札成功 → トースト表示 → 入札履歴更新
3. 価格開示受信 → 価格更新アニメーション → トースト表示
4. 他者の入札受信 → 入札履歴更新 → 入札ボタン有効化
5. 商品終了受信 → 落札者判定 → トースト表示 → 次の商品選択
6. WebSocket切断 → 再接続 → 接続状態表示

### 14.3 成功基準

**機能要件:**
- [ ] 入札ボタンクリックで入札が実行される
- [ ] ポイント不足の場合ボタンが無効化される
- [ ] 価格開示イベントを受信して価格が更新される
- [ ] 他者の入札を受信して入札履歴が更新される
- [ ] 商品終了時に落札者判定が正しく行われる
- [ ] WebSocket切断時に自動再接続される

**パフォーマンス要件:**
- [ ] 初期ロード時間: 2秒以内
- [ ] 入札レスポンス時間: 500ms以内
- [ ] WebSocketイベント受信から画面更新まで: 100ms以内

**セキュリティ要件:**
- [ ] JWT認証なしでAPIアクセス不可（401）
- [ ] 入札者ロールでない場合アクセス不可（403）
- [ ] 同時入札で1人のみ成功（分散ロック）
- [ ] ポイント残高がマイナスにならない（CHECKコンストラインション）

---

## 15. 環境変数

既存の環境変数を使用します。新規追加なし。

**バックエンド（.env）:**
```
JWT_SECRET=your-secret-key-here
REDIS_URL=redis:6379
DATABASE_URL=postgres://auction_user:auction_pass@postgres:5432/auction_db
```

**フロントエンド（.env）:**
```
VITE_API_BASE_URL=http://localhost/api
VITE_WS_URL=ws://localhost/ws
```

---

## 16. 実装手順

### フェーズ1: バックエンド基盤（推定: 4時間）

- [x] Domain層: bidder_points.goの作成（構造体定義）
- [x] Repository層: bid_repository.goの作成（CRUD操作）
- [x] Repository層: point_repository.goの作成（ポイント更新、履歴作成）
- [x] Service層: point_service.goの作成（ポイント取得、キャッシュ）
- [x] Handler層: bid_handler.goの作成（GetPoints実装）
- [x] ルート定義: GET /api/bidder/pointsの追加
- [ ] 動作確認: curlでポイント取得APIをテスト

### フェーズ2: バックエンド入札機能（推定: 6時間）

- [ ] Service層: bid_service.goの作成（PlaceBid実装）
  - [ ] バリデーション
  - [ ] Redis分散ロック
  - [ ] トランザクション処理
  - [ ] ポイント更新
  - [ ] 入札レコード作成
  - [ ] Redis Pub/Sub発行
- [ ] Handler層: PlaceBid、GetBidHistoryの実装
- [ ] ルート定義: POST /api/bidder/items/:id/bid、GET /api/bidder/items/:id/bidsの追加
- [ ] 動作確認: curlで入札APIをテスト
- [ ] エラーハンドリングの確認（ポイント不足、価格変更、同時入札）

### フェーズ3: フロントエンド基盤（推定: 3時間）

- [ ] API Service: bidderBidApi.jsの作成（getPoints, placeBid, getBidHistory）
- [ ] Pinia Store: bidderAuctionLive.jsの作成（state, computed, actions）
- [ ] ルート定義: /bidder/auctions/:id/liveの追加
- [ ] 動作確認: ルートにアクセスしてコンソールエラーなし

### フェーズ4: フロントエンドUI（推定: 8時間）

- [ ] コンポーネント: PointsDisplay.vueの作成（ポイント表示）
- [ ] コンポーネント: BidPanel.vueの作成（入札パネル）
- [ ] コンポーネント: BidderBidHistory.vueの作成（入札履歴）
- [ ] コンポーネント: ItemTabs.vueの作成（商品タブ）
- [ ] メインビュー: BidderAuctionLive.vueの作成
  - [ ] レイアウト実装（PC 3カラム、モバイル縦並び）
  - [ ] 初期化処理（onMounted）
  - [ ] WebSocket接続
  - [ ] エラーハンドリング
- [ ] 動作確認: ブラウザで画面表示、データ取得確認

### フェーズ5: リアルタイム機能（推定: 5時間）

- [ ] WebSocketイベントハンドラー実装（Piniaストア）
  - [ ] onPriceOpened: 価格開示受信
  - [ ] onBidPlaced: 入札受信
  - [ ] onItemEnded: 商品終了受信
  - [ ] onAuctionEnded: オークション終了受信
- [ ] UIアニメーション実装
  - [ ] 価格更新アニメーション（animate-price-update）
  - [ ] トースト通知
  - [ ] ローディングスピナー
- [ ] 入札ボタンの状態管理実装
- [ ] 動作確認: 主催者画面と入札者画面を並べてリアルタイム更新確認

### フェーズ6: テスト＆バグ修正（推定: 4時間）

- [ ] バックエンドユニットテスト作成
- [ ] フロントエンドコンポーネントテスト作成
- [ ] E2Eテスト実施（手動）
  - [ ] 入札フロー
  - [ ] 価格開示受信
  - [ ] 他者の入札受信
  - [ ] 商品終了受信
  - [ ] WebSocket再接続
- [ ] バグ修正
- [ ] パフォーマンステスト

**合計推定時間: 30時間**

---

## 17. 成功基準

### 17.1 機能要件

- [ ] 入札者がオークションに参加し、リアルタイムで価格開示を確認できる
- [ ] 入札者が現在価格で入札できる（ボタン1クリック）
- [ ] ポイント残高がリアルタイムで更新される（利用可能、予約済み、合計）
- [ ] 入札履歴が自分と他者で区別して表示される
- [ ] 価格開示イベントを受信して画面が即座に更新される
- [ ] 他者の入札を受信して入札履歴が更新され、再入札可能になる
- [ ] 商品終了時に落札者判定が正しく行われ、ポイントが消費/返却される
- [ ] WebSocket切断時に自動再接続が動作する（最大5回）
- [ ] 入札ボタンの有効/無効が条件に応じて正しく切り替わる
- [ ] トースト通知が成功/エラー時に適切に表示される

### 17.2 非機能要件

- [ ] レスポンシブデザイン: モバイル/タブレット/PCで正しく表示される
- [ ] アクセシビリティ: キーボード操作、スクリーンリーダー対応
- [ ] パフォーマンス: 初期ロード2秒以内、入札レスポンス500ms以内
- [ ] セキュリティ: JWT認証、分散ロック、ポイント整合性保証
- [ ] エラーハンドリング: すべてのエラーシナリオで適切なメッセージ表示

---

## 18. 次のステップ

この画面の実装完了後、以下の機能を実装します:

1. **オークション結果画面（入札者用）**: 落札した商品、落札できなかった商品の確認
2. **ポイント履歴画面（入札者用）**: ポイントの付与、使用、返却の履歴表示
3. **通知機能**: 価格開示、入札、落札をメールやプッシュ通知で通知
4. **商品のお気に入り機能**: 気になる商品をウォッチリストに追加
5. **入札アラート**: 自分の入札が他者に上書きされた時のアラート

---

## 重要な実装ファイル

以下のファイルが実装の中心となります:

**バックエンド:**
- `backend/internal/handler/bid_handler.go` - 入札APIハンドラー（新規作成）
- `backend/internal/service/bid_service.go` - 入札ビジネスロジック（新規作成）
- `backend/internal/service/point_service.go` - ポイント管理（新規作成）
- `backend/internal/repository/bid_repository.go` - 入札データアクセス（新規作成）
- `backend/internal/repository/point_repository.go` - ポイントデータアクセス（新規作成）
- `backend/internal/domain/bidder_points.go` - ポイントドメインモデル（新規作成）
- `backend/cmd/api/main.go` - ルート定義追加

**フロントエンド:**
- `frontend/src/views/bidder/BidderAuctionLive.vue` - メイン画面（新規作成）
- `frontend/src/stores/bidderAuctionLive.js` - Piniaストア（新規作成）
- `frontend/src/services/bidderBidApi.js` - 入札API（新規作成）
- `frontend/src/components/bidder/BidPanel.vue` - 入札パネル（新規作成）
- `frontend/src/components/bidder/PointsDisplay.vue` - ポイント表示（新規作成）
- `frontend/src/components/bidder/BidderBidHistory.vue` - 入札履歴（新規作成）
- `frontend/src/router/index.js` - ルート追加

---

**このプランに基づいて実装を開始してください。不明点があれば随時質問してください。**
