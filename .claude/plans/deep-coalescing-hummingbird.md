# Admin Dashboard Implementation Plan

## Overview

Implement the admin dashboard screen (`/admin/dashboard`) for the real-time auction system. This is the main landing page for administrators (system_admin and auctioneer roles) after login, providing an at-a-glance view of system statistics, recent activities, and quick action buttons.

**Current State**: The dashboard route exists with only placeholder content showing user information.

**Goal**: Create a fully functional dashboard with:
- Real-time system statistics (4 cards)
- Recent activity feeds (latest bids, new bidders, ended auctions)
- Role-based quick action buttons
- Responsive design for all devices

## Requirements Summary

### System Statistics Cards
1. **Active Auctions Count** - Number of currently active auctions
2. **Today's Bids Count** - Total bids placed today
3. **Registered Bidders Count** - Total active bidders in system
4. **Points Circulation** - Total points distributed across all bidders

### Recent Activities
1. **Latest Bids** (last 5) - Auction name, bidder name, price, timestamp
2. **New Bidders** (last 5) - Email, display name, registration date (system_admin only)
3. **Ended Auctions** (last 5) - Title, winner, final price, end time

### Quick Actions (Role-Based)
1. **Create New Auction** - Available to all admins
2. **Create New Bidder** - system_admin only
3. **Grant Points** - system_admin only

## Implementation Approach

### Phase 1: Backend API Implementation

**New Files to Create:**
- `backend/internal/domain/dashboard.go` - Dashboard data structures
- `backend/internal/repository/dashboard_repository.go` - Database queries for statistics
- `backend/internal/service/dashboard_service.go` - Business logic for dashboard data
- `backend/internal/handler/dashboard_handler.go` - HTTP handlers for dashboard endpoints

**Existing Files to Modify:**
- `backend/cmd/api/main.go` - Add dashboard routes to admin protected group

**API Endpoints:**
1. `GET /api/admin/dashboard/stats` - Returns all 4 statistics
2. `GET /api/admin/dashboard/activities` - Returns recent activities (role-based filtering)

**Database Queries:**
- Active auctions: `SELECT COUNT(*) FROM auctions WHERE status = 'active'`
- Today's bids: `SELECT COUNT(*) FROM bids WHERE bid_at >= CURRENT_DATE`
- Total bidders: `SELECT COUNT(*) FROM bidders WHERE status = 'active'`
- Points circulation: `SELECT SUM(total_points) FROM bidder_points`
- Recent bids: Join `bids`, `items`, `bidders` tables, order by `bid_at DESC`, limit 5
- New bidders: Select from `bidders`, order by `created_at DESC`, limit 5
- Ended auctions: Join `auctions`, `items`, `bidders`, filter by `status = 'ended'`, limit 5

**Role-Based Logic:**
- `new_bidders` data only returned for system_admin role
- Auctioneer receives empty array or field is omitted

### Phase 2: Frontend State Management

**New Files to Create:**
- `frontend/src/services/api/dashboardApi.js` - API client for dashboard endpoints
- `frontend/src/stores/dashboardStore.js` - Pinia store for dashboard state
- `frontend/src/utils/timeFormatter.js` - Utility for relative time formatting (e.g., "3分前")

**Pinia Store Structure:**
```javascript
{
  state: {
    stats: { activeAuctions, todayBids, totalBidders, totalPoints },
    activities: { recentBids, newBidders, endedAuctions },
    loading: false,
    error: null
  },
  actions: {
    fetchStats(),
    fetchActivities(),
    fetchAll()
  }
}
```

### Phase 3: Frontend UI Components

**New Components to Create:**
- `frontend/src/components/admin/StatsCard.vue` - Reusable statistics card with icon, title, and value
- `frontend/src/components/admin/RecentBidsList.vue` - List of recent bids with relative timestamps
- `frontend/src/components/admin/NewBiddersList.vue` - List of new bidders (system_admin only)
- `frontend/src/components/admin/EndedAuctionsList.vue` - List of recently ended auctions
- `frontend/src/components/admin/QuickActions.vue` - Role-based action buttons

**Existing Files to Modify:**
- `frontend/src/views/admin/DashboardView.vue` - Replace placeholder with complete dashboard layout

**Component Layout:**
```
DashboardView.vue
├── Header ("ダッシュボード")
├── Stats Grid (4 StatsCard components)
├── Content Grid (3 columns on desktop, 1-2 on mobile/tablet)
│   ├── RecentBidsList
│   ├── EndedAuctionsList
│   └── Sidebar
│       ├── QuickActions
│       └── NewBiddersList (if system_admin)
```

### Phase 4: Styling and Responsiveness

**Responsive Breakpoints:**
- **Mobile (< 768px)**: 1 column layout, stats cards stacked vertically
- **Tablet (768px - 1279px)**: 2 column grid for stats, mixed layout for content
- **Desktop (≥ 1280px)**: 4 column grid for stats, 3 column content layout

**Design System:**
- Use Shadcn Vue components (Card, Button) with Tailwind CSS
- Color scheme: Blue for primary actions, gray for backgrounds
- Hover effects on stats cards (shadow elevation)
- Loading skeleton screens while fetching data
- Error states with retry button

### Phase 5: Testing and Validation

**Backend Tests:**
- Repository layer: Test each statistics query returns correct count
- Service layer: Test role-based data filtering
- Handler layer: Test JWT auth, role verification, response format

**Frontend Tests:**
- Component tests: Verify stats cards render with correct props
- Store tests: Verify API calls and state updates
- Integration tests: Verify role-based UI rendering

## Critical Files

### Must Create (Backend):
1. `backend/internal/handler/dashboard_handler.go` - Core API logic
2. `backend/internal/repository/dashboard_repository.go` - Database aggregations
3. `backend/internal/service/dashboard_service.go` - Business logic
4. `backend/internal/domain/dashboard.go` - Data structures

### Must Create (Frontend):
1. `frontend/src/stores/dashboardStore.js` - State management
2. `frontend/src/components/admin/StatsCard.vue` - Statistics display
3. `frontend/src/services/api/dashboardApi.js` - API client
4. `frontend/src/utils/timeFormatter.js` - Time formatting utility

### Must Modify:
1. `backend/cmd/api/main.go` - Add dashboard routes
2. `frontend/src/views/admin/DashboardView.vue` - Complete dashboard UI

## Implementation Steps

### Step 1: Backend Foundation (3-4 hours)
1. Create domain structs for dashboard data (`dashboard.go`)
2. Implement repository methods for each statistic (`dashboard_repository.go`)
3. Implement service layer with role-based filtering (`dashboard_service.go`)
4. Create HTTP handlers for 2 endpoints (`dashboard_handler.go`)
5. Register routes in `main.go` under admin protected group
6. Test with curl/Postman

### Step 2: Frontend State Layer (2-3 hours)
1. Create API client functions (`dashboardApi.js`)
2. Create Pinia store with state, getters, actions (`dashboardStore.js`)
3. Create time formatter utility (`timeFormatter.js`)
4. Test API integration

### Step 3: Frontend UI Components (4-5 hours)
1. Build `StatsCard.vue` - Icon, title, large number display
2. Build `RecentBidsList.vue` - Table/list with relative timestamps
3. Build `NewBiddersList.vue` - Conditional rendering based on role
4. Build `EndedAuctionsList.vue` - Winners and final prices
5. Build `QuickActions.vue` - Role-based button visibility
6. Update `DashboardView.vue` - Integrate all components, add grid layout

### Step 4: Styling and Polish (2-3 hours)
1. Implement responsive breakpoints (mobile/tablet/desktop)
2. Add loading skeletons during data fetch
3. Add error handling UI with retry button
4. Add hover effects and transitions
5. Verify accessibility (keyboard nav, screen readers)

### Step 5: Testing (2-3 hours)
1. Write backend unit tests for repository/service/handler
2. Write frontend component tests
3. Manual testing across roles (system_admin vs auctioneer)
4. Manual testing across devices (responsive)

**Total Estimated Time**: 13-18 hours

## Success Criteria

- [ ] All 4 statistics display correctly on dashboard load
- [ ] Recent activities show up to 5 items each
- [ ] New bidders list only visible to system_admin
- [ ] Quick action buttons work and respect role permissions
- [ ] Loading states show during data fetch
- [ ] Error states display with retry option
- [ ] Responsive design works on mobile/tablet/desktop
- [ ] All links navigate to correct pages
- [ ] Data updates when navigating back to dashboard
- [ ] All tests pass

## Security Considerations

1. **Authentication**: All endpoints require valid JWT token
2. **Authorization**: Dashboard accessible to system_admin and auctioneer only (bidders get 403)
3. **Data Filtering**: New bidders data filtered out for auctioneer role
4. **No PII Leakage**: Only display names shown, no sensitive bidder UUIDs exposed
5. **Rate Limiting**: Consider adding rate limits to statistics endpoints (future)

## Performance Considerations

1. **Database Queries**: All use COUNT/SUM aggregations with LIMIT for efficiency
2. **Caching Strategy**: Statistics could be cached for 30-60 seconds (future enhancement)
3. **Lazy Loading**: Components load data on mount, not all upfront
4. **Pagination**: Activity lists limited to 5 items to keep response small

## Future Enhancements

1. **Auto-refresh**: WebSocket or polling for real-time updates every 30 seconds
2. **Charts/Graphs**: Visualize bid trends, auction history over time
3. **Date Filters**: Allow filtering statistics by date range
4. **Export**: CSV/PDF reports of statistics
5. **Customization**: Allow admins to choose which stats to display
6. **Notifications**: Badge indicators for new activities

## References

- Requirements: `docs/screen_list.md` (lines 35-59)
- Database Schema: `docs/database_definition.md`
- Existing Patterns: `frontend/src/views/admin/AuctionListView.vue`
- Auth Store: `frontend/src/stores/auth.js` (has `isSystemAdmin`, `isAuctioneer`)
- API Routes: `backend/cmd/api/main.go`
