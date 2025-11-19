# ============================================
# Multi-stage Dockerfile for Production
# ============================================

# ビルドステージ
FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係のインストール
COPY package*.json ./
RUN npm ci

# ソースコードのコピー
COPY . .

# 本番ビルド
RUN npm run build

# 実行ステージ
FROM nginx:alpine

# ビルドされた静的ファイルをコピー
COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx設定（SPAルーティング対応）
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
