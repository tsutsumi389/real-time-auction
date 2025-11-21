# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Real-time auction system modeled after horse sales, featuring point-based bidding with auctioneer-driven price disclosure. Built with Go backend (REST API + WebSocket servers), Vue.js 3 frontend, PostgreSQL, and Redis.

**Core Architectural Principles:**
- **Dual Server Architecture**: Separate Go servers for REST API ([cmd/api](backend/cmd/api)) and WebSocket ([cmd/ws](backend/cmd/ws))
- **Auctioneer-Driven Pricing**: Auctioneer discloses prices; bidders can only bid at disclosed prices (no free-form bidding)
- **Three-Role System**: `system_admin` (full access), `auctioneer` (auction management), `bidder` (bidding only)
- **Point-Based Economy**: Virtual points only, no real money transactions

## Development Commands

All commands run from project root using the Makefile:

**Essential Commands:**
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
```

**Access Points:**
- Frontend: http://localhost
- REST API: http://localhost/api
- WebSocket: ws://localhost/ws

## Backend Architecture (Go)

### Directory Structure
```
backend/
  cmd/
    api/main.go        # REST API server (port 8080)
    ws/main.go         # WebSocket server (port 8081)
  internal/            # NOT YET IMPLEMENTED - planned structure:
    domain/           # Domain models (structs)
    repository/       # Data access (PostgreSQL via GORM, Redis via go-redis)
    service/          # Business logic (bid validation, point management)
    handler/          # HTTP handlers (Gin)
    ws/              # WebSocket hub, client, event handlers
    middleware/       # JWT auth, CORS, logging
  pkg/                # NOT YET IMPLEMENTED - shared utilities
  migrations/         # SQL migrations (001-004 implemented)
```

### Technology Stack
- **Gin**: HTTP router and middleware
- **Gorilla WebSocket**: RFC 6455-compliant WebSocket (planned for hub pattern)
- **GORM**: PostgreSQL ORM
- **go-redis**: Redis client (Pub/Sub, caching, sessions)
- **golang-jwt/jwt**: JWT authentication
- **Air**: Hot reload (.air.toml for API, .air.ws.toml for WebSocket)

### Database Schema (PostgreSQL)

Migrations in [backend/migrations/](backend/migrations/) define:

**User Tables:**
- `bidders`: UUID primary key (for privacy), email, password_hash, display_name, status
- `admins`: BIGSERIAL primary key, email, password_hash, role (system_admin/auctioneer), status
- `bidder_points`: 1:1 with bidders, tracks total/available/reserved points

**Auction Tables:**
- `auctions`: id, title, status (pending/active/ended/cancelled), starting_price, current_price, winner_id
- `items`: auction items with JSONB metadata field (flexible for horses, art, vehicles, etc.)
- `item_media`: 1:N with items, supports images/videos with display_order
- `bids`: auction_id, bidder_id, price, bid_at, is_winning flag
- `price_history`: tracks auctioneer's price disclosures, had_bid flag

**Audit Tables:**
- `point_history`: Complete audit trail of all point operations (grant/reserve/release/consume/refund) with before/after balances

**Key Constraints:**
- `CHECK (available_points + reserved_points <= total_points)` on bidder_points
- Triggers auto-update `updated_at`, create bidder_points on bidder insert, update is_winning flags
- Email validation via regex CHECK constraints

### Critical Business Logic (Planned Implementation)

**Bidding Flow:**
1. Client sends WebSocket event: `{"type":"auction:bid", "auction_id":1, "price":150}`
2. Handler validates: JWT token, bidder role, price matches current disclosed price
3. Service checks: sufficient available_points, no concurrent bids (Redis lock)
4. Transaction: INSERT bid, UPDATE bidder_points (available → reserved)
5. Redis: Update auction state, Pub/Sub broadcast to all WebSocket servers
6. Hub goroutines broadcast to connected clients in auction room

**Price Disclosure (Auctioneer):**
1. POST /api/auctions/:id/open-price with new price
2. Check if previous price received bids (has_bid flag in Redis)
3. If no bids on previous price → finalize auction with last bidder as winner
4. If bids exist → update current_price, INSERT price_history, broadcast via WebSocket

**Point Management:**
- Grant (system_admin): +total_points, +available_points
- Reserve (on bid): -available_points, +reserved_points
- Consume (on win): -reserved_points
- Release (outbid): +available_points, -reserved_points
- All operations recorded in point_history with full before/after state

## Frontend (Vue.js 3)

### Structure
```
frontend/src/
  views/         # Page components
  components/    # Reusable components (minimal implementation)
  router/        # Vue Router
  stores/        # Pinia state management (not implemented)
  services/      # Axios API clients (not implemented)
```

### Tech Stack
- Vue 3 Composition API (`<script setup>` style)
- Vite (dev server with HMR)
- **Shadcn Vue + Tailwind CSS**: Modern UI component library for design system
- Pinia, Axios (installed but not used yet)
- Native WebSocket API

### UI Design Framework
**Shadcn Vue + Tailwind CSS** is used for modern, customizable UI components:
- **Tailwind CSS**: Utility-first CSS framework for rapid styling
- **Shadcn Vue**: Unstyled, accessible component primitives (based on Radix Vue)
- **Class Variance Authority**: Type-safe component variants
- **Lucide Vue Next**: Icon library (optional)

**Key Benefits:**
- Components copied directly into project (no external dependency bloat)
- Full customization control over styling and behavior
- Built-in dark mode support
- Accessibility (ARIA) compliant
- Smooth animations perfect for real-time updates (bid notifications, price changes)

**Design Philosophy:**
- Modern, clean aesthetic (Vercel/Linear/GitHub style)
- Glassmorphism and micro-interactions for premium feel
- Responsive design for desktop and mobile browsers
- Consistent design tokens via Tailwind config

### Environment Variables
- `VITE_API_BASE_URL`: REST API base (default: http://localhost/api)
- `VITE_WS_URL`: WebSocket URL (default: ws://localhost/ws)

## Infrastructure

### Nginx Reverse Proxy
Routes defined in [nginx/nginx.conf](nginx/nginx.conf):
- `/api/*` → api:8080 (60s timeout, CORS enabled)
- `/ws` → ws:8081 (WebSocket upgrade, 7-day timeout, buffering off)
- `/*` → frontend:5173 (Vite dev server with HMR)

### Redis Usage (Planned)
- Sessions: `session:{token}`
- Auction state: `auction:{id}:current_price`, `auction:{id}:status`, `auction:{id}:has_bid`, `auction:{id}:last_bidder`
- Pub/Sub channels: `auction:started`, `auction:bid`, `auction:price_open`, `auction:ended`
- Distributed locks for bid concurrency control

### Docker Compose
Services: postgres, redis, api, ws, frontend, nginx
- PostgreSQL: 5432 (user: auction_user, db: auction_db)
- Redis: 6379
- Health checks ensure DB readiness before backend starts
- Volume mounts enable hot reload for api, ws, frontend

## Critical Implementation Notes

### WebSocket Server Design (Not Yet Implemented)
Planned hub pattern:
```go
// ws/hub.go
type Hub struct {
    clients    map[*Client]bool
    rooms      map[int64][]*Client  // auction_id -> clients
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}
// Goroutines: readPump(), writePump() per client
// Ping/Pong: 30s interval, 60s timeout
```

### Role-Based Access Control
JWT claims must include `role` field. Validate on every endpoint/WebSocket event:
- `system_admin`: User management, point grants, system-wide access
- `auctioneer`: Create/start/end auctions, register items, disclose prices
- `bidder`: Bid on active auctions, view own points/history

### Common Development Pitfalls
1. **Environment Variables**: `.env` auto-created by `make up`, but change `JWT_SECRET` for production
2. **Docker Networking**: Use `postgres:5432` not `localhost` inside containers
3. **WebSocket Upgrade**: Nginx must set `Upgrade` header (already configured)
4. **CORS**: Currently allows `*` in dev; restrict in production via `CORS_ORIGINS`
5. **Goroutine Cleanup**: Ensure readPump/writePump exit cleanly on WebSocket disconnect

### Testing & Debugging
```bash
# Health checks
curl http://localhost/api/health
curl http://localhost/ws/health

# Database inspection
make shell-postgres
\dt                          # List tables
\d+ bidders                  # Table schema
SELECT * FROM active_auctions_view;

# Redis inspection
make shell-redis
KEYS *
GET auction:1:current_price
```

## Implementation Planning Guidelines

When creating implementation plans for new features or screens, follow these guidelines:

### Plan Document Format

**Location**: Store all implementation plans in `docs/plan/` directory

**File Naming**: Use descriptive names with underscores (e.g., `admin_login_implementation_plan.md`, `auction_create_implementation_plan.md`)

**Target Audience**: Write plans for non-technical stakeholders (project managers, designers, business analysts)
- Use plain Japanese with technical terms explained
- Avoid including actual code examples
- Focus on "what" needs to be built, not "how" to code it
- Use diagrams, tables, and flowcharts instead of code

### Required Sections

1. **Overview** (概要)
   - Purpose and goals
   - Target users
   - Key features

2. **Screen Layout** (画面レイアウト)
   - ASCII diagrams showing screen composition
   - Design requirements (colors, fonts, spacing)
   - Responsive design breakpoints

3. **Input Fields and Validation** (入力項目とバリデーション)
   - Table format with: label, input type, placeholder, required/optional, constraints
   - Validation rules with user-friendly error messages
   - Validation timing (onBlur, onSubmit, realtime)

4. **Processing Flow** (処理フロー)
   - Step-by-step flow diagrams
   - Success/failure scenarios
   - Error handling with user-facing messages

5. **Security Requirements** (セキュリティ要件)
   - Data handling (encryption, hashing)
   - Token management (JWT, sessions)
   - Common vulnerabilities addressed (SQL injection, XSS, CSRF)

6. **Database Design** (データベース設計)
   - Table structures in table format
   - Constraints and relationships
   - Test data specifications

7. **API Specification** (API仕様)
   - Endpoint paths and HTTP methods
   - Request/response JSON examples
   - Error response catalog

8. **Backend Implementation Requirements** (バックエンド実装要件)
   - Layer responsibilities (Domain, Repository, Service, Handler)
   - Describe "what each layer does" not "how to code it"
   - Middleware requirements
   - Technology stack

9. **Frontend Implementation Requirements** (フロントエンド実装要件)
   - Directory structure
   - Component descriptions (what they do)
   - State management requirements
   - Routing configuration

10. **Screen Transitions** (画面遷移)
    - Transition flow diagrams
    - Navigation guard behavior
    - Redirect rules

11. **UI Behavior Specification** (UIの動作仕様)
    - User interactions (click, keyboard, focus)
    - Loading states
    - Success/error feedback

12. **Responsive Design** (レスポンシブデザイン)
    - Mobile, tablet, desktop breakpoints
    - Layout adjustments per device

13. **Accessibility** (アクセシビリティ)
    - Keyboard navigation
    - Screen reader support
    - Focus management

14. **Test Requirements** (テスト要件)
    - Test scenarios (not test code)
    - Success criteria
    - Edge cases to cover

15. **Environment Variables** (環境変数)
    - Required configuration
    - Example values

16. **Implementation Steps** (実装手順)
    - Phase-by-phase breakdown with checkboxes (`- [ ]` format)
    - Time estimates per phase
    - Dependencies between phases
    - Total estimated time

17. **Success Criteria** (成功基準)
    - Checklist of completion requirements

18. **Next Steps** (次のステップ)
    - Related features to implement next

**Example**: See [docs/plan/admin_login_implementation_plan.md](docs/plan/admin_login_implementation_plan.md) for reference format

### What NOT to Include in Plans

- ❌ Actual Go/TypeScript/Vue code examples
- ❌ Implementation details (variable names, function signatures)
- ❌ Code snippets in any programming language (except JSON for API examples)
- ❌ Technical jargon without explanation

### What TO Include in Plans

- ✅ Plain language descriptions
- ✅ Diagrams and flowcharts
- ✅ Tables and structured data
- ✅ User-facing error messages
- ✅ Business logic flow
- ✅ JSON examples for API requests/responses
- ✅ SQL table definitions in table format
- ✅ Technology choices with rationale

## Additional Documentation

- [Database Definition](docs/database_definition.md): Full schema, triggers, views, and optimization strategies
- [Architecture Document](docs/architecture.md): System design and technical decisions
- [Screen List](docs/screen_list.md): Planned UI screens and user flows
- [Copilot Instructions](.github/copilot-instructions.md): Detailed development workflow (also applicable to Claude Code)
- [Implementation Plans](docs/plan/): Feature-specific implementation plans (non-technical format)

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
