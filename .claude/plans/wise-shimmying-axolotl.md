# 入札者側ラグジュアリーデザイン実装計画

## 概要
入札者向け画面を競走馬取引にふさわしいラグジュアリーなデザインに刷新する。管理者側のデザインは現状維持。ダークモード対応あり。

---

## カラーパレット

### ライトモード
| 用途 | 色名 | HSL値 | HEX |
|------|------|-------|-----|
| Primary (Gold) | シャンパンゴールド | 43 74% 49% | #D4A84B |
| Secondary (Burgundy) | ボルドー | 345 65% 35% | #8B2942 |
| Accent (Green) | レーシンググリーン | 152 45% 35% | #3A7D5C |
| Background | クリーム | 40 33% 98% | #FAF8F5 |
| Card | ウォームホワイト | 40 20% 99% | #FDFCFB |
| Muted | ウォームグレー | 40 15% 93% | #EDEBE8 |
| Platinum | プラチナ | 220 15% 75% | #B8BCC6 |

### ダークモード
| 用途 | HSL値 |
|------|-------|
| Background | 220 25% 8% |
| Card | 220 22% 12% |
| Primary (Gold) | 43 70% 55% |
| Secondary (Burgundy) | 345 55% 45% |

---

## タイポグラフィ

- **見出し**: Playfair Display (Google Fonts)
- **本文**: Inter (Google Fonts)
- **価格表示**: Inter tabular-nums

```html
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Playfair+Display:wght@400;500;600;700&display=swap" rel="stylesheet">
```

---

## 修正ファイル一覧

### 1. 設定ファイル

#### `frontend/index.html`
- Google Fontsのリンク追加

#### `frontend/tailwind.config.js`
- fontFamily追加 (sans: Inter, serif: Playfair Display)
- auction色パレット拡張 (gold-light, gold-dark, burgundy系, green系, platinum, silver, cream, ivory)
- status色追加
- boxShadow追加 (luxury, luxury-lg, luxury-xl, gold-glow, inner-luxury)
- keyframes追加 (gold-shimmer, pulse-gold, winning-glow)
- animation追加
- backgroundImage追加 (gold-gradient, burgundy-gradient, luxury-gradient)

#### `frontend/src/assets/main.css`
- :root CSS変数の追加・更新
- .dark CSS変数の追加・更新
- .bidder-theme クラスの追加（入札者向けテーマ適用用）
- body基本スタイルの更新（フォント設定）
- 見出し用スタイル (.heading, h1, h2, h3)

### 2. 入札者向けコンポーネント

#### `frontend/src/components/bidder/BidderHeader.vue`
- 背景色: bg-background → ゴールドボーダー下線追加
- ロゴテキスト: セリフフォント、バーガンディ色
- ナビリンク: ホバー時ゴールド、下線アニメーション
- アバター: グリーン→バーガンディグラデーション、ゴールドボーダー
- ログイン/登録ボタン: ゴールドアクセント

#### `frontend/src/views/bidder/BidderLoginView.vue`
- 背景: blue-50/indigo-100 → cream/muted グラデーション
- カード: ゴールド上部ボーダー、luxury-lg シャドウ
- 見出し: セリフフォント、バーガンディ色
- 入力フィールド: ゴールドフォーカスリング
- ボタン: ゴールド背景

#### `frontend/src/components/bidder/AuctionCard.vue`
- カード: ゴールド上部ボーダー、platinum/50 ボーダー、luxury シャドウ
- タイトル: セリフフォント、バーガンディ色
- アイコン: platinum色
- ボタン: ゴールド背景、グリーン「参加する」ボタン

#### `frontend/src/components/bidder/BidPanel.vue`
- 価格表示エリア: cream背景、ゴールドボーダー、inner-luxury シャドウ
- 「Current Price」ラベル: ゴールド色、tracking-widest
- 価格: セリフフォント
- 入札ボタン: ゴールドグラデーション、gold-glow シャドウ
- 最高入札者表示: ゴールドボーダー、pulse-gold アニメーション

#### `frontend/src/components/bidder/ItemDetailModal.vue`
- バックドロップ: black/60 → burgundy-dark/80
- モーダル: ゴールドボーダー、luxury-xl シャドウ
- LOTバッジ: ゴールド背景
- タイトル: セリフフォント、バーガンディ色
- 価格セクション: cream背景、ゴールドアクセント
- 閉じるボタン: ホバー時ゴールド

#### `frontend/src/components/bidder/AuctionStatusBadge.vue`
- pending: amber系 (変更なし)
- active: auction-green系 + パルスドット
- ended: muted系
- cancelled: red系

### 3. 共有UIコンポーネント（入札者向けバリアント追加）

#### `frontend/src/components/ui/Button.vue`
- 新バリアント追加: `luxury` (ゴールドグラデーション)
- 新バリアント追加: `luxury-outline` (ゴールドボーダー)
- 新バリアント追加: `luxury-secondary` (バーガンディ)

### 4. レイアウト

#### `frontend/src/layouts/DefaultLayout.vue`
- bidder-themeクラスを適用するラッパー追加

---

## 実装順序

1. **Phase 1: 基盤設定**
   - index.html (Google Fonts)
   - tailwind.config.js (色、フォント、シャドウ、アニメーション)
   - main.css (CSS変数、基本スタイル)

2. **Phase 2: 共有コンポーネント**
   - Button.vue (luxuryバリアント追加)
   - DefaultLayout.vue (テーマクラス適用)

3. **Phase 3: 入札者コンポーネント**
   - BidderHeader.vue
   - AuctionStatusBadge.vue
   - AuctionCard.vue
   - BidPanel.vue
   - ItemDetailModal.vue

4. **Phase 4: 画面**
   - BidderLoginView.vue

5. **Phase 5: 確認・調整**
   - ダークモード動作確認
   - レスポンシブ確認

---

## 注意事項

- 管理者側コンポーネント (`/components/admin/`, `/views/admin/`) は変更しない
- 既存の機能・アクセシビリティは維持
- `prefers-reduced-motion` 対応を維持
- 既存のアニメーション (ripple, price-pop等) は保持しつつラグジュアリー調整
