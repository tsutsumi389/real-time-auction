# ============================================
# Multi-stage Dockerfile for Production
# ============================================

# ビルドステージ
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 依存関係のインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# バイナリのビルド（API用とWS用）
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o ws ./cmd/ws

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# ビルドされたバイナリをコピー
COPY --from=builder /build/api .
COPY --from=builder /build/ws .

# タイムゾーン設定
ENV TZ=Asia/Tokyo

# デフォルトでAPIサーバーを起動
# docker-compose.ymlでコマンドをオーバーライド可能
CMD ["./api"]
