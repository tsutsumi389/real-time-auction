# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Real-time auction system modeled after horse sales, featuring point-based bidding with auctioneer-driven price disclosure. Built with Go backend (REST API + WebSocket servers), Vue.js 3 frontend, PostgreSQL, and Redis.

**Core Architectural Principles:**
- **Dual Server Architecture**: Separate Go servers for REST API and WebSocket
- **Auctioneer-Driven Pricing**: Auctioneer discloses prices; bidders can only bid at disclosed prices
- **Three-Role System**: `system_admin` (full access), `auctioneer` (auction management), `bidder` (bidding only)
- **Point-Based Economy**: Virtual points only, no real money transactions

## Development Commands

All commands run from project root using the Makefile:

```bash
# Start all services (auto-creates .env on first run)
make up

# View logs (all services or specific: make logs service=api)
make logs

# Stop services / Clean everything including volumes
make down
make clean

# Access service shells
make shell-api          # REST API server container
make shell-ws           # WebSocket server container
make shell-postgres     # PostgreSQL (psql)
make shell-redis        # Redis CLI

# Database migrations (using golang-migrate)
make db-migrate         # Apply migrations
make db-migrate-down    # Rollback one migration
make db-status          # Check migration status
make db-create-migration name=description  # Create new migration

# Build commands (must be run inside containers)
docker compose exec api go build -o /app/bin/api ./cmd/api        # Build REST API server
docker compose exec ws go build -o /app/bin/ws ./cmd/ws           # Build WebSocket server
docker compose exec frontend npm run build                         # Build frontend (production)
```

**Access Points:**
- Frontend: http://localhost
- REST API: http://localhost/api
- WebSocket: ws://localhost/ws

## Technology Stack

**Backend (Go):**
- Gin (HTTP router), Gorilla WebSocket, GORM (PostgreSQL), go-redis, golang-jwt/jwt, Air (hot reload)

**Frontend (Vue.js 3):**
- Vue 3 Composition API, Vite, Shadcn Vue + Tailwind CSS, Pinia, Axios, Native WebSocket API

**Infrastructure:**
- PostgreSQL (5432), Redis (6379), Nginx (reverse proxy)

## Database Schema

See [docs/database_definition.md](docs/database_definition.md) for full schema details.

**Key Tables:**
- `bidders`, `admins`: User management with separate tables for security
- `auctions`, `items`, `item_media`: Auction and item management
- `bids`, `price_history`: Bidding and price disclosure tracking
- `bidder_points`, `point_history`: Point management and audit trail

## Development Guidelines

**Backend Guidelines**: See [docs/rule/backend.md](docs/rule/backend.md)
- Directory structure, technology stack, critical business logic, WebSocket design, RBAC, Redis usage

**Frontend Guidelines**: See [docs/rule/frontend.md](docs/rule/frontend.md)
- API-first development approach, technology stack, UI design framework, directory structure

**Implementation Planning**: See [docs/rule/planning.md](docs/rule/planning.md)
- Plan document format, required sections, target audience, what to include/exclude

## Testing & Debugging

```bash
# Health checks
curl http://localhost/api/health
curl http://localhost/ws/health

# Database inspection
make shell-postgres
\dt                          # List tables
\d+ bidders                  # Table schema

# Redis inspection
make shell-redis
KEYS *
GET auction:1:current_price
```

## Additional Documentation

- [Database Definition](docs/database_definition.md): Full schema, triggers, views, and optimization strategies
- [Architecture Document](docs/architecture.md): System design and technical decisions
- [Screen List](docs/screen_list.md): Planned UI screens and user flows
- [Copilot Instructions](.github/copilot-instructions.md): Detailed development workflow
- [Implementation Plans](docs/plan/): Feature-specific implementation plans (non-technical format)
- [Development Rules](docs/rule/): Backend, frontend, and planning guidelines

## Key Design Decisions

**Why Separate Bidders and Admins Tables?**
- Security: Physical separation reduces privilege escalation risks
- Optimization: Different data structures (bidders need points, admins need roles)
- Scalability: Bidders table will grow to 10k+ records, admins ~100

**Why UUID for Bidders, BIGSERIAL for Admins?**
- Bidders: Privacy (no sequence inference), distributed ID generation, URL guessing prevention
- Admins: Internal only, better performance, simpler debugging

**Why JSONB for Item Metadata?**
- Flexibility: Support diverse auction categories (horses, art, vehicles) without schema changes
- Extensibility: Add new fields without migrations
- GIN indexing: Efficient queries on JSON attributes
