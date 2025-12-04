# Backend Development Guidelines

## Directory Structure

```
backend/
  cmd/
    api/main.go        # REST API server (port 8080)
    ws/main.go         # WebSocket server (port 8081)
  internal/
    domain/           # Domain models (structs)
    repository/       # Data access (PostgreSQL via GORM, Redis via go-redis)
    service/          # Business logic (bid validation, point management)
    handler/          # HTTP handlers (Gin)
    ws/              # WebSocket hub, client, event handlers
    middleware/       # JWT auth, CORS, logging
  pkg/                # Shared utilities
  migrations/         # SQL migrations
```

## Technology Stack

- **Gin**: HTTP router and middleware
- **Gorilla WebSocket**: RFC 6455-compliant WebSocket
- **GORM**: PostgreSQL ORM
- **go-redis**: Redis client (Pub/Sub, caching, sessions)
- **golang-jwt/jwt**: JWT authentication
- **Air**: Hot reload (.air.toml for API, .air.ws.toml for WebSocket)

## Critical Business Logic

### Bidding Flow

1. Client sends WebSocket event: `{"type":"auction:bid", "auction_id":1, "price":150}`
2. Handler validates: JWT token, bidder role, price matches current disclosed price
3. Service checks: sufficient available_points, no concurrent bids (Redis lock)
4. Transaction: INSERT bid, UPDATE bidder_points (available → reserved)
5. Redis: Update auction state, Pub/Sub broadcast to all WebSocket servers
6. Hub goroutines broadcast to connected clients in auction room

### Price Disclosure (Auctioneer)

1. POST /api/auctions/:id/open-price with new price
2. Check if previous price received bids (has_bid flag in Redis)
3. If no bids on previous price → finalize auction with last bidder as winner
4. If bids exist → update current_price, INSERT price_history, broadcast via WebSocket

### Point Management

- Grant (system_admin): +total_points, +available_points
- Reserve (on bid): -available_points, +reserved_points
- Consume (on win): -reserved_points
- Release (outbid): +available_points, -reserved_points
- All operations recorded in point_history with full before/after state

## WebSocket Server Design

Hub pattern implementation:

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

## Role-Based Access Control

JWT claims must include `role` field. Validate on every endpoint/WebSocket event:
- `system_admin`: User management, point grants, system-wide access
- `auctioneer`: Create/start/end auctions, register items, disclose prices
- `bidder`: Bid on active auctions, view own points/history

## Common Development Pitfalls

1. **Environment Variables**: `.env` auto-created by `make up`, but change `JWT_SECRET` for production
2. **Docker Networking**: Use `postgres:5432` not `localhost` inside containers
3. **WebSocket Upgrade**: Nginx must set `Upgrade` header (already configured)
4. **CORS**: Currently allows `*` in dev; restrict in production via `CORS_ORIGINS`
5. **Goroutine Cleanup**: Ensure readPump/writePump exit cleanly on WebSocket disconnect

## Redis Usage

- Sessions: `session:{token}`
- Auction state: `auction:{id}:current_price`, `auction:{id}:status`, `auction:{id}:has_bid`, `auction:{id}:last_bidder`
- Pub/Sub channels: `auction:started`, `auction:bid`, `auction:price_open`, `auction:ended`
- Distributed locks for bid concurrency control
