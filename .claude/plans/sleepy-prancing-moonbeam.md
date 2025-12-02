# BidderAuctionListView UI/UX 改善実装計画

## 概要

BidderAuctionListView とその関連コンポーネント全体の UI/UX を改善します。視覚デザインの向上、ユーザビリティの改善、アクセシビリティの強化を目的とした包括的な改善を行います。

### 改善の重点領域

1. **視覚デザインの向上**: 洗練されたカードデザイン、スケルトンローディング、滑らかなアニメーション
2. **ユーザビリティの向上**: フィルタUIの改善（native select → カスタムSelect）
3. **新規コンポーネント導入**: Select/Dropdown、Skeleton
4. **アクセシビリティ強化**: ARIA属性、キーボードナビゲーション

---

## 実装フェーズ

### フェーズ1: 基盤コンポーネントの実装（最優先）

#### 1.1 Skeleton コンポーネントの実装

**新規ファイル**: `frontend/src/components/ui/Skeleton.vue`

```vue
目的: ローディング状態を視覚的に表現する汎用スケルトンコンポーネント
機能:
- Pulse/Shimmer アニメーション
- 複数サイズ対応（sm, md, lg, full）
- カスタム幅・高さ対応
- 角丸対応（rounded variant）
```

**Props**:
- `variant`: 'default' | 'circular' | 'rectangular'
- `width`: string (例: '100%', '200px')
- `height`: string (例: '20px', '100px')
- `rounded`: boolean

**実装詳細**:
- CVA（Class Variance Authority）でバリアント管理
- Tailwind `animate-pulse` を基盤にカスタムアニメーション追加
- ARIA属性: `aria-busy="true"`, `aria-label="読み込み中"`

#### 1.2 AuctionCardSkeleton コンポーネントの実装

**新規ファイル**: `frontend/src/components/bidder/AuctionCardSkeleton.vue`

```vue
目的: AuctionCard 専用のスケルトン表示
構成:
- 画像エリア（高さ 192px）
- タイトル（2行）
- 説明文（2行）
- メタ情報（3行）
- ボタンエリア（2個）
```

**実装詳細**:
- Skeleton コンポーネントを組み合わせて構成
- AuctionCard と同じ構造・サイズを維持
- レスポンシブ対応（sm, lg ブレークポイント）

#### 1.3 Tailwind Config の拡張

**修正ファイル**: `frontend/tailwind.config.js`

追加する keyframes と animations:
```javascript
keyframes: {
  // 既存のaccordion系に追加
  "shimmer": {
    "0%": { backgroundPosition: "-1000px 0" },
    "100%": { backgroundPosition: "1000px 0" }
  },
  "fade-in": {
    from: { opacity: "0" },
    to: { opacity: "1" }
  },
  "slide-in-up": {
    from: { transform: "translateY(10px)", opacity: "0" },
    to: { transform: "translateY(0)", opacity: "1" }
  },
  "scale-in": {
    from: { transform: "scale(0.95)", opacity: "0" },
    to: { transform: "scale(1)", opacity: "1" }
  }
},
animation: {
  // 既存のaccordion系に追加
  "shimmer": "shimmer 2s infinite linear",
  "fade-in": "fade-in 0.3s ease-out",
  "slide-in-up": "slide-in-up 0.4s ease-out",
  "scale-in": "scale-in 0.3s ease-out"
}
```

---

### フェーズ2: Select コンポーネント群の実装

#### 2.1 Select コンポーネント実装（9ファイル）

Radix Vue の Select プリミティブを活用した、アクセシブルなカスタムセレクトコンポーネント群を実装します。

**新規ファイル群**:

1. **`frontend/src/components/ui/Select.vue`**
   - ルートコンポーネント、context provider

2. **`frontend/src/components/ui/SelectTrigger.vue`**
   - トリガーボタン（現在の選択値を表示）
   - ChevronDown アイコン付き
   - フォーカスリング対応

3. **`frontend/src/components/ui/SelectContent.vue`**
   - ドロップダウンコンテンツコンテナ
   - Teleport で body に描画
   - スクロール対応
   - Enter/Leave アニメーション

4. **`frontend/src/components/ui/SelectItem.vue`**
   - 個別選択肢
   - ホバー・選択状態の視覚フィードバック
   - Check アイコン表示

5. **`frontend/src/components/ui/SelectValue.vue`**
   - 選択値表示コンポーネント
   - プレースホルダー対応

6. **`frontend/src/components/ui/SelectGroup.vue`**
   - 選択肢グループコンテナ

7. **`frontend/src/components/ui/SelectLabel.vue`**
   - グループラベル

8. **`frontend/src/components/ui/SelectSeparator.vue`**
   - 区切り線

9. **`frontend/src/components/ui/SelectScrollButton.vue`**
   - スクロールボタン（上下）

**実装のポイント**:
- Radix Vue の `@radix-vue/select` パッケージ活用
- provide/inject パターンで状態共有
- キーボードナビゲーション完全対応（↑↓キー、Enter、Escape）
- ARIA属性自動付与（role, aria-selected, aria-expanded等）
- Shadcn Vue デザインシステムに準拠したスタイリング

#### 2.2 Lucide Icons のインストール確認

**確認事項**: `lucide-vue-next` パッケージがインストール済み

使用するアイコン:
- `ChevronDown`: Select トリガー
- `Check`: 選択済みアイテム
- `Search`: 検索バー
- `Filter`: フィルタ
- `ArrowUpDown`: ソート
- `X`: クリアボタン
- `Loader2`: ローディングスピナー

---

### フェーズ3: 既存コンポーネントの改善

#### 3.1 AuctionCard.vue の改善

**修正ファイル**: `frontend/src/components/bidder/AuctionCard.vue`

**視覚デザインの改善**:

1. **カードのホバーエフェクト強化**:
   ```vue
   変更前: hover:shadow-lg
   変更後: hover:shadow-xl hover:scale-[1.02] transition-all duration-300
   ```

2. **画像オーバーレイの追加**:
   - グラデーションオーバーレイで画像とバッジのコントラスト向上
   - `bg-gradient-to-t from-black/30 to-transparent`

3. **Lucide アイコンの導入**:
   - 既存のSVGアイコンを Lucide Icons に置き換え
   - `Package2`: 出品物数
   - `Calendar`: 開始日時
   - `Clock`: 更新日時
   - `ImageOff`: 画像なし状態

4. **ボタンを Shadcn Button コンポーネントに置き換え**:
   ```vue
   <Button variant="default" size="default" @click="handleViewDetails">
     詳細を見る
   </Button>
   <Button v-if="auction.status === 'active'" variant="default" size="default" @click="handleJoinAuction">
     参加する
   </Button>
   ```

5. **タイポグラフィとスペーシングの最適化**:
   - タイトル: `text-xl font-bold` → `text-xl font-semibold`
   - 説明: `text-gray-600` → `text-muted-foreground`
   - メタ情報のギャップ: `space-y-2` → `space-y-2.5`

**アクセシビリティの改善**:
- ボタンに `aria-label` 追加
- カード全体に `role="article"` 追加
- 画像に適切な `alt` 属性（既存実装済み）

#### 3.2 AuctionCardGrid.vue の改善

**修正ファイル**: `frontend/src/components/bidder/AuctionCardGrid.vue`

**主な変更**:

1. **スケルトンローディングの実装**:
   ```vue
   <template v-if="loading">
     <AuctionCardSkeleton v-for="i in skeletonCount" :key="`skeleton-${i}`" />
   </template>
   ```
   - `skeletonCount` prop 追加（デフォルト: 6）
   - 初回ローディング時にスケルトンカードを表示

2. **TransitionGroup によるアニメーション**:
   ```vue
   <TransitionGroup
     name="card-list"
     tag="div"
     class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
   >
     <AuctionCard
       v-for="auction in auctions"
       :key="auction.id"
       :auction="auction"
       @join-auction="handleJoinAuction"
     />
   </TransitionGroup>
   ```

3. **Stagger アニメーションのスタイル追加**:
   ```css
   .card-list-enter-active {
     transition: all 0.4s ease-out;
   }
   .card-list-enter-from {
     opacity: 0;
     transform: translateY(20px);
   }
   .card-list-move {
     transition: transform 0.4s ease;
   }
   ```

4. **ARIA Live Region の追加**:
   ```vue
   <div aria-live="polite" aria-atomic="true" class="sr-only">
     {{ auctions.length }}件のオークションを表示中
   </div>
   ```

#### 3.3 AuctionFilters.vue の改善

**修正ファイル**: `frontend/src/components/bidder/AuctionFilters.vue`

**主な変更**:

1. **Native Select を Shadcn Select に置き換え**:
   ```vue
   変更前:
   <select v-model="selectedSort" @change="handleSortChange">
     <option v-for="option in sortOptions" :value="option.value">
       {{ option.label }}
     </option>
   </select>

   変更後:
   <Select v-model="selectedSort" @update:modelValue="handleSortChange">
     <SelectTrigger class="w-full sm:w-[240px]">
       <SelectValue placeholder="並び替え" />
     </SelectTrigger>
     <SelectContent>
       <SelectItem v-for="option in sortOptions" :value="option.value">
         {{ option.label }}
       </SelectItem>
     </SelectContent>
   </Select>
   ```

2. **Lucide Icons の追加**:
   - `Filter`: フィルタセクションラベル
   - `ArrowUpDown`: ソートラベル

3. **ステータスフィルタのアクセシビリティ向上**:
   ```vue
   <button
     v-for="option in statusOptions"
     :key="option.value"
     @click="handleStatusChange(option.value)"
     :aria-pressed="currentStatus === option.value"
     :aria-label="`${option.label}のオークションを表示`"
     role="radio"
   >
   ```

4. **レイアウトの改善**:
   - フレックスレイアウトの最適化
   - モバイルでの可読性向上

#### 3.4 AuctionSearchBar.vue の改善

**修正ファイル**: `frontend/src/components/bidder/AuctionSearchBar.vue`

**主な変更**:

1. **Native Input を Shadcn Input に置き換え**:
   ```vue
   <Input
     v-model="localKeyword"
     type="text"
     :placeholder="placeholder"
     :disabled="loading"
     @keyup.enter="handleSearch"
     class="flex-1"
   />
   ```

2. **Lucide Icons の追加**:
   - `Search`: 検索アイコン（Input の先頭）
   - `Loader2`: ローディング中アイコン（回転アニメーション）
   - `X`: クリアボタン

3. **Shadcn Button の導入**:
   ```vue
   <Button
     @click="handleSearch"
     :disabled="loading"
     variant="default"
   >
     <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
     <Search v-else class="mr-2 h-4 w-4" />
     {{ searchButtonText }}
   </Button>
   ```

4. **アクセシビリティの改善**:
   ```vue
   <div role="search" aria-label="オークション検索">
     <Input aria-label="検索キーワード" />
   </div>
   ```

#### 3.5 LoadingSpinner.vue の改善

**修正ファイル**: `frontend/src/components/ui/LoadingSpinner.vue`

**主な変更**:

1. **CVA によるバリアント管理の導入**:
   ```javascript
   const spinnerVariants = cva(
     "inline-block animate-spin rounded-full border-2 border-solid border-current border-r-transparent",
     {
       variants: {
         size: {
           xs: "h-3 w-3",
           sm: "h-4 w-4",
           md: "h-6 w-6",
           lg: "h-12 w-12",
           xl: "h-16 w-16",
           "2xl": "h-20 w-20"
         },
         color: {
           primary: "text-primary",
           secondary: "text-secondary",
           white: "text-white",
           muted: "text-muted-foreground"
         }
       },
       defaultVariants: {
         size: "md",
         color: "primary"
       }
     }
   )
   ```

2. **新しい Props の追加**:
   - `color`: 'primary' | 'secondary' | 'white' | 'muted'
   - `size`: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'

3. **アクセシビリティの強化**:
   ```vue
   <div
     :class="spinnerClasses"
     role="status"
     aria-busy="true"
     :aria-label="text || '読み込み中'"
   >
     <span class="sr-only">{{ text || '読み込み中' }}</span>
   </div>
   ```

#### 3.6 BidderAuctionListView.vue の改善

**修正ファイル**: `frontend/src/views/BidderAuctionListView.vue`

**主な変更**:

1. **AuctionCardSkeleton のインポート追加**:
   ```javascript
   import AuctionCardSkeleton from '@/components/bidder/AuctionCardSkeleton.vue'
   ```

2. **無限スクロールトリガーの改善**:
   - `setTimeout` を `nextTick()` に変更
   ```javascript
   import { nextTick } from 'vue'

   watch(() => auctionStore.pagination.hasMore, async (hasMore) => {
     if (hasMore) {
       await nextTick()
       setupIntersectionObserver()
     }
   })
   ```

3. **エラー表示の改善**:
   ```vue
   <Alert v-if="auctionStore.error" variant="destructive" class="mb-6">
     <AlertCircle class="h-4 w-4" />
     <AlertTitle>エラー</AlertTitle>
     <AlertDescription>{{ auctionStore.error }}</AlertDescription>
   </Alert>
   ```

4. **handleViewDetails の実装**:
   ```javascript
   const handleViewDetails = (auction) => {
     router.push({
       name: 'bidder-auction-detail',
       params: { id: auction.id }
     })
   }
   ```

---

## アクセシビリティチェックリスト

### ARIA 属性
- [ ] ボタンに `aria-label` または `aria-labelledby`
- [ ] Select に `role="combobox"`, `aria-expanded`, `aria-controls`
- [ ] フィルタボタンに `role="radio"`, `aria-pressed`
- [ ] ローディング状態に `aria-busy="true"`, `role="status"`
- [ ] Live region に `aria-live="polite"`, `aria-atomic="true"`

### キーボードナビゲーション
- [ ] Tab キーでフォーカス移動
- [ ] Enter キーで検索実行
- [ ] Escape キーで Select を閉じる
- [ ] ↑↓キーで Select 内を移動
- [ ] Space キーでフィルタボタンをトグル

### スクリーンリーダー対応
- [ ] 画像に適切な `alt` 属性
- [ ] 視覚的な情報を `sr-only` クラスでテキスト化
- [ ] フォーム要素に `<label>` 関連付け
- [ ] 状態変更を live region で通知

---

## 実装ファイル一覧

### 新規作成ファイル（11個）

1. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/Skeleton.vue`
2. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/bidder/AuctionCardSkeleton.vue`
3. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/Select.vue`
4. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectTrigger.vue`
5. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectContent.vue`
6. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectItem.vue`
7. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectValue.vue`
8. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectGroup.vue`
9. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectLabel.vue`
10. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectSeparator.vue`
11. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/SelectScrollButton.vue`

### 修正ファイル（7個）

1. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/tailwind.config.js`
2. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/bidder/AuctionCard.vue`
3. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/bidder/AuctionCardGrid.vue`
4. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/bidder/AuctionFilters.vue`
5. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/bidder/AuctionSearchBar.vue`
6. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/components/ui/LoadingSpinner.vue`
7. `/Users/tsutsumi/devlop/web/real-time-auction/frontend/src/views/BidderAuctionListView.vue`

---

## 実装スケジュール

### フェーズ1: 基盤コンポーネント（2-3日）
- [ ] Skeleton.vue 実装
- [ ] AuctionCardSkeleton.vue 実装
- [ ] Tailwind Config 拡張
- [ ] AuctionCardGrid.vue にスケルトン統合
- [ ] AuctionCard.vue デザイン改善
- [ ] カードアニメーション実装

### フェーズ2: Select コンポーネント群（2-3日）
- [ ] Select コンポーネント群（9ファイル）実装
- [ ] Lucide Icons 統合
- [ ] AuctionFilters.vue 改善
- [ ] AuctionSearchBar.vue 改善
- [ ] キーボードナビゲーションテスト
- [ ] アクセシビリティテスト

### フェーズ3: 仕上げと最適化（1-2日）
- [ ] LoadingSpinner.vue 改善
- [ ] BidderAuctionListView.vue 最終調整
- [ ] レスポンシブデザイン確認
- [ ] ブラウザ互換性テスト
- [ ] パフォーマンス最適化
- [ ] ドキュメント更新

**総所要時間**: 5-8営業日

---

## 成功基準

### 視覚デザイン
- [ ] スケルトンローディングが滑らかに表示される
- [ ] カードホバー時のアニメーションが自然
- [ ] カード出現時の Stagger アニメーションが機能
- [ ] デザイントークンが一貫して適用されている

### ユーザビリティ
- [ ] カスタム Select が直感的に操作できる
- [ ] 検索とフィルタリングがスムーズ
- [ ] ローディング状態が明確に伝わる
- [ ] エラー状態が分かりやすい

### アクセシビリティ
- [ ] キーボードのみで全操作が可能
- [ ] スクリーンリーダーで内容を理解できる
- [ ] WCAG 2.1 AA レベルに準拠
- [ ] フォーカスインジケーターが明確

### パフォーマンス
- [ ] 初回レンダリングが高速
- [ ] アニメーションが 60fps で動作
- [ ] 無限スクロールがスムーズ
- [ ] バンドルサイズが適切

---

## 技術的注意事項

### Radix Vue の活用
- `@radix-vue/select` パッケージを使用
- provide/inject パターンで状態管理
- 自動的な ARIA 属性付与を活用

### CVA（Class Variance Authority）
- タイプセーフなコンポーネントバリアント管理
- `cva()` 関数でスタイル定義
- `cn()` ヘルパーでクラス結合

### Lucide Icons
- Tree-shakable なアイコンライブラリ
- 必要なアイコンのみインポート
- 一貫したサイズとスタイル

### パフォーマンス最適化
- TransitionGroup の `mode` 属性を適切に設定
- Intersection Observer の `rootMargin` 調整
- スケルトン表示数の最適化（6個）

---

## 備考

この計画は、既存の実装を最大限活用しながら、段階的かつ安全に UI/UX を改善するよう設計されています。各フェーズは独立して実装可能で、途中でのレビューや調整も容易です。

実装中に発見される課題や改善点は、適宜この計画に反映させながら進めることを推奨します。
