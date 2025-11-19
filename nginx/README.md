# Nginx Configuration

Nginxリバースプロキシ設定

## 概要

開発環境用のNginx設定です。REST API、WebSocket、フロントエンドへのルーティングを提供します。

## ルーティング

```
http://localhost/
├── /api/*        → REST APIサーバー (api:8080)
├── /ws/*         → WebSocketサーバー (ws:8081)
└── /*            → フロントエンド (frontend:5173)
```

## 主な機能

### REST APIプロキシ（/api/*）
- リバースプロキシ設定
- CORS対応（開発環境用に全オリジン許可）
- プリフライトリクエスト（OPTIONS）対応
- タイムアウト: 60秒

### WebSocketプロキシ（/ws/*）
- WebSocketアップグレード対応
- 長時間接続サポート（7日間）
- バッファリング無効化（リアルタイム通信最適化）
- Sticky Session対応準備

### フロントエンドプロキシ（/）
- Vite開発サーバーへのプロキシ
- HMR（Hot Module Replacement）サポート
- WebSocketアップグレード対応

### ヘルスチェック
- `/api/health` - APIサーバーのヘルスチェック
- `/ws/health` - WebSocketサーバーのヘルスチェック

## 本番環境との違い

開発環境の設定は以下の点で本番環境と異なります：

1. **CORS設定**: 全オリジン許可（本番では特定オリジンのみ）
2. **SSL/TLS**: HTTP のみ（本番では HTTPS 必須）
3. **レート制限**: 未設定（本番では IP/ユーザーごとに制限）
4. **静的ファイル**: Vite dev server 経由（本番では直接配信）
5. **キャッシュ**: 無効（本番では適切なキャッシュ設定）

## 本番環境設定（TODO）

本番環境では以下の追加設定が必要です：

- SSL/TLS証明書の設定
- CORS オリジン制限
- レート制限（nginx-limit-req-module）
- セキュリティヘッダー（CSP、X-Frame-Options等）
- 静的ファイルの直接配信
- Gzip圧縮の最適化
- アクセスログの適切な管理
