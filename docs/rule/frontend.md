# Frontend Development Guidelines

## API-First Development

**Before implementing any frontend feature:**
1. **Check existing API endpoints** in [backend/internal/handler/](../../backend/internal/handler/) to understand available routes
2. **Review API contracts** - Check request/response structures in handler implementations
3. **Verify middleware** - Check authentication and authorization requirements in [backend/internal/middleware/](../../backend/internal/middleware/)
4. **Test endpoints** - Use curl or similar tools to verify API behavior before frontend integration
5. **Never assume API structure** - Always verify the actual implementation, don't rely on plans or assumptions

**API Documentation Locations:**
- Handler implementations: [backend/internal/handler/](../../backend/internal/handler/)
- Route definitions: [backend/cmd/api/main.go](../../backend/cmd/api/main.go)
- Request/Response models: [backend/internal/domain/](../../backend/internal/domain/)
- Implementation plans (reference only): [docs/plan/](../plan/)

## Technology Stack

- **Vue 3 Composition API** (`<script setup>` style)
- **Vite**: Dev server with HMR
- **Shadcn Vue + Tailwind CSS**: Modern UI component library for design system
- **Pinia**: State management (installed but not fully implemented)
- **Axios**: HTTP client (installed but not fully implemented)
- **Native WebSocket API**: Real-time communication

## UI Design Framework

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

## Directory Structure

```
frontend/src/
  views/         # Page components
  components/    # Reusable components
  router/        # Vue Router configuration
  stores/        # Pinia state management
  services/      # Axios API clients
```

## Environment Variables

- `VITE_API_BASE_URL`: REST API base (default: http://localhost/api)
- `VITE_WS_URL`: WebSocket URL (default: ws://localhost/ws)

## Code Style

- Use Composition API with `<script setup>` syntax
- Prefer composables for reusable logic
- Use TypeScript for type safety
- Follow Vue 3 style guide and best practices
