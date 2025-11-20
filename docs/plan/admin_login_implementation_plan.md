# 管理者ログイン画面（1.1.1）実装プラン

**作成日**: 2025年11月20日
**対象画面**: 管理者ログイン画面（Admin Login）
**パス**: `/admin/login`
**権限**: 未認証ユーザー

---

## 1. 概要

管理者（system_admin / auctioneer）が認証してシステムにアクセスするためのログイン画面。JWT認証を使用し、成功時にはトークンをクライアント側に保存してダッシュボードへリダイレクトする。

---

## 2. 実装範囲

### 2.1 バックエンド実装（Go）

#### 2.1.1 データモデル（Domain）
**ファイル**: `backend/internal/domain/admin.go`

```go
// Admin represents an administrator or auctioneer user
type Admin struct {
    ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // Never expose in JSON
    DisplayName  string    `gorm:"type:varchar(100)" json:"display_name"`
    Role         string    `gorm:"type:varchar(20);not null;check:role IN ('system_admin','auctioneer')" json:"role"`
    Status       string    `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active','suspended','deleted')" json:"status"`
    CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Admin) TableName() string {
    return "admins"
}
```

**JWT Claims構造体**:
```go
// JWTClaims represents JWT token claims
type JWTClaims struct {
    UserID      int64  `json:"user_id"`
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
    Role        string `json:"role"`
    UserType    string `json:"user_type"` // "admin" or "bidder"
    jwt.RegisteredClaims
}
```

---

#### 2.1.2 Repository層（Data Access）
**ファイル**: `backend/internal/repository/admin_repository.go`

```go
type AdminRepository interface {
    FindByEmail(ctx context.Context, email string) (*domain.Admin, error)
    FindByID(ctx context.Context, id int64) (*domain.Admin, error)
}

type adminRepository struct {
    db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
    return &adminRepository{db: db}
}

// FindByEmail retrieves an admin by email address
func (r *adminRepository) FindByEmail(ctx context.Context, email string) (*domain.Admin, error) {
    var admin domain.Admin
    if err := r.db.WithContext(ctx).
        Where("email = ? AND status != ?", email, "deleted").
        First(&admin).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrAdminNotFound
        }
        return nil, err
    }
    return &admin, nil
}

// FindByID retrieves an admin by ID
func (r *adminRepository) FindByID(ctx context.Context, id int64) (*domain.Admin, error) {
    var admin domain.Admin
    if err := r.db.WithContext(ctx).
        Where("id = ? AND status != ?", id, "deleted").
        First(&admin).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrAdminNotFound
        }
        return nil, err
    }
    return &admin, nil
}
```

**カスタムエラー定義**:
```go
var (
    ErrAdminNotFound = errors.New("admin not found")
)
```

---

#### 2.1.3 Service層（Business Logic）
**ファイル**: `backend/internal/service/auth_service.go`

```go
type AuthService interface {
    AdminLogin(ctx context.Context, email, password string) (string, *domain.Admin, error)
    ValidateJWT(tokenString string) (*domain.JWTClaims, error)
}

type authService struct {
    adminRepo repository.AdminRepository
    jwtSecret string
    jwtExpiry time.Duration
}

func NewAuthService(adminRepo repository.AdminRepository, jwtSecret string) AuthService {
    return &authService{
        adminRepo: adminRepo,
        jwtSecret: jwtSecret,
        jwtExpiry: 24 * time.Hour, // 24 hours
    }
}

// AdminLogin authenticates an admin and returns JWT token
func (s *authService) AdminLogin(ctx context.Context, email, password string) (string, *domain.Admin, error) {
    // Find admin by email
    admin, err := s.adminRepo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, repository.ErrAdminNotFound) {
            return "", nil, ErrInvalidCredentials
        }
        return "", nil, err
    }

    // Check if admin is suspended
    if admin.Status == "suspended" {
        return "", nil, ErrAccountSuspended
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
        return "", nil, ErrInvalidCredentials
    }

    // Generate JWT token
    token, err := s.generateJWT(admin)
    if err != nil {
        return "", nil, err
    }

    return token, admin, nil
}

// generateJWT creates a new JWT token for admin
func (s *authService) generateJWT(admin *domain.Admin) (string, error) {
    claims := &domain.JWTClaims{
        UserID:      admin.ID,
        Email:       admin.Email,
        DisplayName: admin.DisplayName,
        Role:        admin.Role,
        UserType:    "admin",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "auction-system",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

// ValidateJWT validates and parses JWT token
func (s *authService) ValidateJWT(tokenString string) (*domain.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*domain.JWTClaims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}
```

**カスタムエラー定義**:
```go
var (
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrAccountSuspended   = errors.New("account is suspended")
    ErrInvalidToken       = errors.New("invalid token")
)
```

---

#### 2.1.4 Handler層（HTTP Handlers）
**ファイル**: `backend/internal/handler/auth_handler.go`

```go
type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

// LoginRequest represents the login request body
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// LoginResponse represents the login response
type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}

type UserResponse struct {
    ID          int64  `json:"id"`
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
    Role        string `json:"role"`
    UserType    string `json:"user_type"`
}

// AdminLogin handles POST /api/auth/admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "details": err.Error(),
        })
        return
    }

    // Authenticate admin
    token, admin, err := h.authService.AdminLogin(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCredentials) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid email or password",
            })
            return
        }
        if errors.Is(err, service.ErrAccountSuspended) {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Account is suspended",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Internal server error",
        })
        return
    }

    // Return success response
    c.JSON(http.StatusOK, LoginResponse{
        Token: token,
        User: UserResponse{
            ID:          admin.ID,
            Email:       admin.Email,
            DisplayName: admin.DisplayName,
            Role:        admin.Role,
            UserType:    "admin",
        },
    })
}
```

---

#### 2.1.5 Middleware（JWT認証）
**ファイル**: `backend/internal/middleware/auth_middleware.go`

```go
// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // Extract Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization header format",
            })
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Validate token
        claims, err := authService.ValidateJWT(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Set user context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        c.Set("user_type", claims.UserType)

        c.Next()
    }
}

// RequireRole validates that the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "User role not found",
            })
            c.Abort()
            return
        }

        roleStr := userRole.(string)
        for _, allowedRole := range roles {
            if roleStr == allowedRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": "Insufficient permissions",
        })
        c.Abort()
    }
}
```

---

#### 2.1.6 ルーティング設定
**ファイル**: `backend/cmd/api/main.go`

```go
func setupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, authService service.AuthService) {
    api := router.Group("/api")
    {
        // Public routes
        auth := api.Group("/auth")
        {
            auth.POST("/admin/login", authHandler.AdminLogin)
            // auth.POST("/bidder/login", authHandler.BidderLogin) // 将来実装
        }

        // Protected routes (require authentication)
        admin := api.Group("/admin")
        admin.Use(middleware.AuthMiddleware(authService))
        admin.Use(middleware.RequireRole("system_admin", "auctioneer"))
        {
            // 将来的にダッシュボードや他のエンドポイントを追加
        }
    }
}
```

---

### 2.2 フロントエンド実装（Vue 3 + Shadcn）

#### 2.2.1 ディレクトリ構造
```
frontend/src/
  views/
    admin/
      AdminLoginView.vue      # 管理者ログイン画面
  components/
    ui/                       # Shadcn Vue components
      button.vue
      input.vue
      label.vue
      card.vue
      alert.vue
  services/
    api/
      authApi.ts              # 認証API client
  stores/
    authStore.ts              # 認証状態管理（Pinia）
  router/
    index.ts                  # ルーティング設定
  utils/
    auth.ts                   # JWT token管理
```

---

#### 2.2.2 API Client（Axios）
**ファイル**: `frontend/src/services/api/authApi.ts`

```typescript
import axios, { AxiosInstance } from 'axios'

const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface AdminLoginRequest {
  email: string
  password: string
}

export interface AdminLoginResponse {
  token: string
  user: {
    id: number
    email: string
    display_name: string
    role: string
    user_type: string
  }
}

export const authApi = {
  adminLogin: async (credentials: AdminLoginRequest): Promise<AdminLoginResponse> => {
    const response = await apiClient.post<AdminLoginResponse>('/auth/admin/login', credentials)
    return response.data
  },
}
```

---

#### 2.2.3 JWT Token管理
**ファイル**: `frontend/src/utils/auth.ts`

```typescript
const TOKEN_KEY = 'auction_admin_token'
const USER_KEY = 'auction_admin_user'

export interface StoredUser {
  id: number
  email: string
  display_name: string
  role: string
  user_type: string
}

export const tokenManager = {
  getToken: (): string | null => {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken: (token: string): void => {
    localStorage.setItem(TOKEN_KEY, token)
  },

  removeToken: (): void => {
    localStorage.removeItem(TOKEN_KEY)
  },

  getUser: (): StoredUser | null => {
    const userStr = localStorage.getItem(USER_KEY)
    return userStr ? JSON.parse(userStr) : null
  },

  setUser: (user: StoredUser): void => {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },

  removeUser: (): void => {
    localStorage.removeItem(USER_KEY)
  },

  clearAuth: (): void => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  },

  isAuthenticated: (): boolean => {
    return !!tokenManager.getToken()
  },
}
```

---

#### 2.2.4 Pinia Store（状態管理）
**ファイル**: `frontend/src/stores/authStore.ts`

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type AdminLoginRequest, type AdminLoginResponse } from '@/services/api/authApi'
import { tokenManager, type StoredUser } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<StoredUser | null>(tokenManager.getUser())
  const token = ref<string | null>(tokenManager.getToken())
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const userRole = computed(() => user.value?.role || null)

  const login = async (credentials: AdminLoginRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response: AdminLoginResponse = await authApi.adminLogin(credentials)

      token.value = response.token
      user.value = response.user

      tokenManager.setToken(response.token)
      tokenManager.setUser(response.user)

      return true
    } catch (err: any) {
      if (err.response?.status === 401) {
        error.value = 'メールアドレスまたはパスワードが正しくありません'
      } else if (err.response?.status === 403) {
        error.value = 'アカウントが停止されています'
      } else {
        error.value = 'ログインに失敗しました。もう一度お試しください。'
      }
      return false
    } finally {
      isLoading.value = false
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    tokenManager.clearAuth()
  }

  return {
    user,
    token,
    isLoading,
    error,
    isAuthenticated,
    userRole,
    login,
    logout,
  }
})
```

---

#### 2.2.5 Vue Component（ログイン画面）
**ファイル**: `frontend/src/views/admin/AdminLoginView.vue`

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const emailError = ref('')
const passwordError = ref('')

const validateEmail = (): boolean => {
  if (!email.value) {
    emailError.value = 'メールアドレスを入力してください'
    return false
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email.value)) {
    emailError.value = 'メールアドレスの形式が正しくありません'
    return false
  }

  emailError.value = ''
  return true
}

const validatePassword = (): boolean => {
  if (!password.value) {
    passwordError.value = 'パスワードを入力してください'
    return false
  }

  if (password.value.length < 8) {
    passwordError.value = 'パスワードは8文字以上で入力してください'
    return false
  }

  passwordError.value = ''
  return true
}

const handleSubmit = async () => {
  const isEmailValid = validateEmail()
  const isPasswordValid = validatePassword()

  if (!isEmailValid || !isPasswordValid) {
    return
  }

  const success = await authStore.login({
    email: email.value,
    password: password.value,
  })

  if (success) {
    router.push('/admin/dashboard')
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 px-4">
    <Card class="w-full max-w-md shadow-xl">
      <CardHeader class="space-y-1">
        <CardTitle class="text-2xl font-bold tracking-tight">管理者ログイン</CardTitle>
        <CardDescription class="text-sm text-muted-foreground">
          システム管理者またはオークション主催者としてログインしてください
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-4">
        <!-- Error Alert -->
        <Alert v-if="authStore.error" variant="destructive" class="mb-4">
          <AlertDescription>{{ authStore.error }}</AlertDescription>
        </Alert>

        <!-- Email Input -->
        <div class="space-y-2">
          <Label for="email">メールアドレス</Label>
          <Input
            id="email"
            v-model="email"
            type="email"
            placeholder="admin@example.com"
            autocomplete="email"
            :disabled="authStore.isLoading"
            :class="{ 'border-red-500': emailError }"
            @blur="validateEmail"
            @keydown.enter="handleSubmit"
          />
          <p v-if="emailError" class="text-sm text-red-500">{{ emailError }}</p>
        </div>

        <!-- Password Input -->
        <div class="space-y-2">
          <Label for="password">パスワード</Label>
          <Input
            id="password"
            v-model="password"
            type="password"
            placeholder="8文字以上"
            autocomplete="current-password"
            :disabled="authStore.isLoading"
            :class="{ 'border-red-500': passwordError }"
            @blur="validatePassword"
            @keydown.enter="handleSubmit"
          />
          <p v-if="passwordError" class="text-sm text-red-500">{{ passwordError }}</p>
        </div>
      </CardContent>

      <CardFooter>
        <Button
          class="w-full"
          :disabled="authStore.isLoading"
          @click="handleSubmit"
        >
          <span v-if="authStore.isLoading">ログイン中...</span>
          <span v-else>ログイン</span>
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>
```

---

#### 2.2.6 ルーティング設定
**ファイル**: `frontend/src/router/index.ts`

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import { tokenManager } from '@/utils/auth'
import AdminLoginView from '@/views/admin/AdminLoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/admin/login',
      name: 'AdminLogin',
      component: AdminLoginView,
      meta: { requiresAuth: false },
    },
    {
      path: '/admin/dashboard',
      name: 'AdminDashboard',
      component: () => import('@/views/admin/AdminDashboardView.vue'),
      meta: { requiresAuth: true, roles: ['system_admin', 'auctioneer'] },
    },
    // 他のルート定義...
  ],
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const isAuthenticated = tokenManager.isAuthenticated()

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/admin/login')
  } else if (to.path === '/admin/login' && isAuthenticated) {
    next('/admin/dashboard')
  } else {
    next()
  }
})

export default router
```

---

## 3. データベース設定

### 3.1 Seed Data（開発用）
**ファイル**: `backend/migrations/005_seed_admin_users.up.sql`

```sql
-- Insert test admin users (password: "password123" hashed with bcrypt)
INSERT INTO admins (email, password_hash, display_name, role, status)
VALUES
  ('admin@example.com', '$2a$10$rZJ3KQ3X2k.9n5K/Z2xJEeW8Uc1R9vQ8Y6Yq7V0xF5Lm4Kf2Y3zL.', 'システム管理者', 'system_admin', 'active'),
  ('auctioneer@example.com', '$2a$10$rZJ3KQ3X2k.9n5K/Z2xJEeW8Uc1R9vQ8Y6Yq7V0xF5Lm4Kf2Y3zL.', '主催者テスト', 'auctioneer', 'active');
```

---

## 4. セキュリティ要件

### 4.1 バックエンドセキュリティ
1. **パスワードハッシュ化**: bcrypt（cost=10）を使用
2. **JWT Secret**: 環境変数`JWT_SECRET`で管理（最低32文字のランダム文字列）
3. **JWT有効期限**: 24時間
4. **レート制限**: ログインエンドポイントに対して1分間に5回まで（将来実装）
5. **HTTPS**: 本番環境では必須
6. **CORS**: 許可されたオリジンのみ（環境変数`CORS_ORIGINS`）
7. **SQLインジェクション対策**: GORMのプレースホルダー使用

### 4.2 フロントエンドセキュリティ
1. **XSS対策**: Vueの自動エスケープ機能
2. **Token保存**: localStorage（sessionStorageも検討可）
3. **HTTPSのみ**: 本番環境でのToken送信
4. **CSRF対策**: JWTトークンを使用（Cookieは使用しない）

---

## 5. エラーハンドリング

### 5.1 バックエンドエラーレスポンス

| HTTPステータス | エラーコード | メッセージ | 説明 |
|--------------|------------|----------|------|
| 400 | `INVALID_REQUEST` | Invalid request body | リクエストボディの形式が不正 |
| 401 | `INVALID_CREDENTIALS` | Invalid email or password | 認証情報が不正 |
| 403 | `ACCOUNT_SUSPENDED` | Account is suspended | アカウントが停止中 |
| 500 | `INTERNAL_ERROR` | Internal server error | サーバー内部エラー |

### 5.2 フロントエンドエラー表示
- **Network Error**: 「サーバーに接続できません」
- **Timeout**: 「接続がタイムアウトしました」
- **401 Unauthorized**: 「メールアドレスまたはパスワードが正しくありません」
- **403 Forbidden**: 「アカウントが停止されています」
- **500 Server Error**: 「ログインに失敗しました。もう一度お試しください。」

---

## 6. バリデーションルール

### 6.1 バックエンドバリデーション
- **Email**: 必須、メールアドレス形式、最大255文字
- **Password**: 必須、最小8文字

### 6.2 フロントエンドバリデーション（リアルタイム）
- **Email**:
  - 必須: 「メールアドレスを入力してください」
  - 形式: 「メールアドレスの形式が正しくありません」
- **Password**:
  - 必須: 「パスワードを入力してください」
  - 最小長: 「パスワードは8文字以上で入力してください」

---

## 7. UI/UXデザイン要件

### 7.1 デザインシステム（Shadcn + Tailwind）
- **カラーテーマ**: Light/Darkモード対応
- **Primary Color**: Slate（管理画面）
- **Accent Color**: Blue（ボタン、リンク）
- **Typography**: Inter font family

### 7.2 レスポンシブデザイン
- **Mobile**: 画面幅100%、パディング16px
- **Tablet/Desktop**: カード幅最大448px（max-w-md）、中央配置

### 7.3 アクセシビリティ
- **Label**: すべての入力フィールドに関連付け
- **Aria-label**: ボタンとリンクに明確なラベル
- **Keyboard Navigation**: Tab/Enterキーでのフォーム送信
- **Focus Styles**: 明確なフォーカス表示

---

## 8. テストケース

### 8.1 バックエンドテスト（Go）
**ファイル**: `backend/internal/handler/auth_handler_test.go`

```go
func TestAdminLogin_Success(t *testing.T) {
    // Test successful login with valid credentials
}

func TestAdminLogin_InvalidCredentials(t *testing.T) {
    // Test login failure with wrong password
}

func TestAdminLogin_AccountSuspended(t *testing.T) {
    // Test login failure with suspended account
}

func TestAdminLogin_InvalidEmail(t *testing.T) {
    // Test validation error with invalid email format
}
```

### 8.2 フロントエンドテスト（Vitest + Vue Test Utils）
**ファイル**: `frontend/src/views/admin/AdminLoginView.spec.ts`

```typescript
describe('AdminLoginView', () => {
  it('renders login form correctly', () => {})
  it('validates email field on blur', () => {})
  it('validates password field on blur', () => {})
  it('displays error message on failed login', () => {})
  it('redirects to dashboard on successful login', () => {})
  it('disables submit button during loading', () => {})
})
```

---

## 9. 実装順序

1. **Phase 1: バックエンド基盤**
   - [ ] データモデル（Admin, JWTClaims）
   - [ ] Repository層（AdminRepository）
   - [ ] Service層（AuthService）
   - [ ] ミドルウェア（AuthMiddleware, RequireRole）

2. **Phase 2: バックエンドAPI**
   - [ ] Handler層（AuthHandler.AdminLogin）
   - [ ] ルーティング設定
   - [ ] Seed data作成（テスト用管理者アカウント）

3. **Phase 3: フロントエンド基盤**
   - [ ] API Client（authApi.ts）
   - [ ] Token管理（auth.ts）
   - [ ] Pinia Store（authStore.ts）

4. **Phase 4: フロントエンドUI**
   - [ ] Shadcn UIコンポーネント設定（Button, Input, Label, Card, Alert）
   - [ ] AdminLoginView.vue実装
   - [ ] ルーティング設定（Navigation Guard）

5. **Phase 5: テストと検証**
   - [ ] バックエンドユニットテスト
   - [ ] フロントエンドコンポーネントテスト
   - [ ] 統合テスト（E2E）
   - [ ] セキュリティ検証

---

## 10. 環境変数

### 10.1 バックエンド（`.env`）
```bash
# JWT Settings
JWT_SECRET=your-secure-random-secret-key-min-32-chars

# Server Settings
API_PORT=8080
CORS_ORIGINS=http://localhost,http://localhost:5173

# Database Settings
DATABASE_URL=postgres://auction_user:auction_pass@postgres:5432/auction_db?sslmode=disable
```

### 10.2 フロントエンド（`.env`）
```bash
VITE_API_BASE_URL=http://localhost/api
```

---

## 11. 注意事項と制約

1. **パスワードの平文保存禁止**: 必ずbcryptでハッシュ化してから保存
2. **JWT Secretの管理**: 環境変数で管理し、Gitにコミットしない
3. **本番環境**: HTTPS必須、JWT Secretは強力なランダム文字列
4. **ログ出力**: パスワードやトークンをログに出力しない
5. **CORS設定**: 本番環境では特定のドメインのみ許可
6. **レート制限**: ブルートフォース攻撃対策として将来実装必須
7. **トークン有効期限**: 24時間、リフレッシュトークンは将来実装

---

## 12. 次のステップ

この実装プラン完了後、以下の機能を実装予定：

1. **ダッシュボード画面（1.1.2）**: 管理者用のホーム画面
2. **入札者ログイン画面（2.1.1）**: 入札者用の認証画面
3. **ログアウト機能**: JWT無効化（Redis blacklist使用）
4. **リフレッシュトークン**: 長期間のセッション維持
5. **パスワードリセット**: メール認証によるパスワード再設定

---

**作成者**: Claude Code
**レビュー状態**: 未レビュー
**最終更新**: 2025年11月20日
