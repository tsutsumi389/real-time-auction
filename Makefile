.DEFAULT_GOAL := help

# ============================================
# 変数定義
# ============================================
COMPOSE := docker-compose
COMPOSE_FILE := docker-compose.yml

# サービス名
SERVICES := postgres redis api ws frontend nginx
SERVICE_API := api
SERVICE_WS := ws
SERVICE_FRONTEND := frontend
SERVICE_NGINX := nginx
SERVICE_POSTGRES := postgres
SERVICE_REDIS := redis

# ============================================
# ヘルプ
# ============================================
.PHONY: help
help: ## ヘルプを表示
	@echo "Real-Time Auction System - 開発コマンド"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "利用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "例:"
	@echo "  make up          # 全サービスを起動"
	@echo "  make logs        # 全サービスのログを表示"
	@echo "  make logs service=api  # APIサーバーのログのみ表示"
	@echo ""

# ============================================
# 環境構築
# ============================================
.PHONY: setup
setup: ## 初回セットアップ (.env作成)
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "✓ .env created. Please edit .env if needed."; \
	else \
		echo ".env already exists. Skipping..."; \
	fi

# ============================================
# Docker操作
# ============================================
.PHONY: up
up: setup ## 全サービスを起動
	$(COMPOSE) up -d
	@echo ""
	@echo "✓ Services are starting..."
	@echo ""
	@echo "Access URLs:"
	@echo "  Frontend:  http://localhost"
	@echo "  API:       http://localhost/api"
	@echo "  WebSocket: ws://localhost/ws"
	@echo ""
	@echo "Run 'make logs' to view logs"

.PHONY: down
down: ## 全サービスを停止
	$(COMPOSE) down
	@echo "✓ All services stopped"

.PHONY: restart
restart: ## 全サービスを再起動
	$(COMPOSE) restart
	@echo "✓ All services restarted"

.PHONY: ps
ps: ## サービスのステータスを表示
	$(COMPOSE) ps

.PHONY: clean
clean: ## 全サービスを停止し、ボリュームも削除
	@echo "⚠️  This will remove all containers, volumes, and data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(COMPOSE) down -v; \
		echo "✓ All services, volumes, and data removed"; \
	else \
		echo "Cancelled."; \
	fi

# ============================================
# ログ確認
# ============================================
.PHONY: logs
logs: ## ログを表示 (service=<name>で特定サービスのみ)
ifdef service
	$(COMPOSE) logs -f $(service)
else
	$(COMPOSE) logs -f
endif

.PHONY: logs-api
logs-api: ## APIサーバーのログを表示
	$(COMPOSE) logs -f $(SERVICE_API)

.PHONY: logs-ws
logs-ws: ## WebSocketサーバーのログを表示
	$(COMPOSE) logs -f $(SERVICE_WS)

.PHONY: logs-frontend
logs-frontend: ## フロントエンドのログを表示
	$(COMPOSE) logs -f $(SERVICE_FRONTEND)

.PHONY: logs-nginx
logs-nginx: ## Nginxのログを表示
	$(COMPOSE) logs -f $(SERVICE_NGINX)

# ============================================
# 個別サービス操作
# ============================================
.PHONY: start
start: ## 特定サービスを起動 (service=<name>)
ifndef service
	@echo "Error: service parameter is required"
	@echo "Usage: make start service=<service_name>"
	@echo "Available services: $(SERVICES)"
	@exit 1
endif
	$(COMPOSE) start $(service)

.PHONY: stop
stop: ## 特定サービスを停止 (service=<name>)
ifndef service
	@echo "Error: service parameter is required"
	@echo "Usage: make stop service=<service_name>"
	@echo "Available services: $(SERVICES)"
	@exit 1
endif
	$(COMPOSE) stop $(service)

.PHONY: restart-service
restart-service: ## 特定サービスを再起動 (service=<name>)
ifndef service
	@echo "Error: service parameter is required"
	@echo "Usage: make restart-service service=<service_name>"
	@echo "Available services: $(SERVICES)"
	@exit 1
endif
	$(COMPOSE) restart $(service)

# ============================================
# コンテナ内操作
# ============================================
.PHONY: shell-api
shell-api: ## APIサーバーのシェルに入る
	$(COMPOSE) exec $(SERVICE_API) /bin/sh

.PHONY: shell-ws
shell-ws: ## WebSocketサーバーのシェルに入る
	$(COMPOSE) exec $(SERVICE_WS) /bin/sh

.PHONY: shell-frontend
shell-frontend: ## フロントエンドのシェルに入る
	$(COMPOSE) exec $(SERVICE_FRONTEND) /bin/sh

.PHONY: shell-postgres
shell-postgres: ## PostgreSQLのシェルに入る
	$(COMPOSE) exec $(SERVICE_POSTGRES) psql -U auction_user -d auction_db

.PHONY: shell-redis
shell-redis: ## Redisのシェルに入る
	$(COMPOSE) exec $(SERVICE_REDIS) redis-cli

# ============================================
# ビルド
# ============================================
.PHONY: build
build: ## 全サービスをビルド
	$(COMPOSE) build

.PHONY: build-api
build-api: ## APIサーバーをビルド
	$(COMPOSE) build $(SERVICE_API)

.PHONY: build-ws
build-ws: ## WebSocketサーバーをビルド
	$(COMPOSE) build $(SERVICE_WS)

.PHONY: build-frontend
build-frontend: ## フロントエンドをビルド
	$(COMPOSE) build $(SERVICE_FRONTEND)

.PHONY: rebuild
rebuild: ## 全サービスを再ビルドして起動 (キャッシュなし)
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d

# ============================================
# データベース操作
# ============================================
.PHONY: db-migrate
db-migrate: ## データベースマイグレーションを実行
	@echo "Running database migrations..."
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/001_create_initial_tables.up.sql
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/002_create_triggers_and_functions.up.sql
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/003_create_views.up.sql
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/004_insert_initial_data.up.sql
	@echo "✓ Database migration completed"

.PHONY: db-migrate-down
db-migrate-down: ## データベースマイグレーションをロールバック
	@echo "⚠️  Rolling back database migrations..."
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/004_insert_initial_data.down.sql; \
		$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/003_create_views.down.sql; \
		$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/002_create_triggers_and_functions.down.sql; \
		$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db < backend/migrations/001_create_initial_tables.down.sql; \
		echo "✓ Database rollback completed"; \
	else \
		echo "Cancelled."; \
	fi

.PHONY: db-status
db-status: ## データベースの状態を確認
	@echo "Database tables:"
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db -c "\dt"
	@echo ""
	@echo "Database views:"
	@$(COMPOSE) exec -T $(SERVICE_POSTGRES) psql -U auction_user -d auction_db -c "\dv"

.PHONY: db-seed
db-seed: ## シードデータを投入 (マイグレーションに含まれています)
	@echo "Seed data is included in migration 004_insert_initial_data"
	@echo "Run 'make db-migrate' to apply seed data"

.PHONY: db-reset
db-reset: ## データベースをリセット (削除して再作成)
	@echo "⚠️  This will delete all database data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		$(COMPOSE) stop $(SERVICE_POSTGRES); \
		docker volume rm real-time-auction_postgres_data 2>/dev/null || true; \
		$(COMPOSE) up -d $(SERVICE_POSTGRES); \
		echo "✓ Database reset complete"; \
	else \
		echo "Cancelled."; \
	fi

# ============================================
# テスト (準備中)
# ============================================
.PHONY: test
test: ## テストを実行 (準備中)
	@echo "Tests not yet implemented"

.PHONY: test-backend
test-backend: ## バックエンドテストを実行 (準備中)
	@echo "Backend tests not yet implemented"

.PHONY: test-frontend
test-frontend: ## フロントエンドテストを実行 (準備中)
	@echo "Frontend tests not yet implemented"

# ============================================
# その他
# ============================================
.PHONY: prune
prune: ## 未使用のDockerリソースを削除
	@echo "Pruning unused Docker resources..."
	docker system prune -f
	docker volume prune -f
	@echo "✓ Docker cleanup complete"
