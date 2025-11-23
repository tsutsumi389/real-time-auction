package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/service"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
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
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("email", claims.Email)

		// For admin users, also store role
		if claims.UserType == domain.UserTypeAdmin {
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

// RequireAdmin middleware ensures the user is an admin
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		if userType != domain.UserTypeAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware ensures the user has a specific admin role
func RequireRole(requiredRole domain.AdminRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		adminRole, ok := role.(domain.AdminRole)
		if !ok || adminRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSystemAdmin middleware ensures the user is a system admin
func RequireSystemAdmin() gin.HandlerFunc {
	return RequireRole(domain.RoleSystemAdmin)
}

// RequireAuctioneer middleware ensures the user is an auctioneer
func RequireAuctioneer() gin.HandlerFunc {
	return RequireRole(domain.RoleAuctioneer)
}

// RequireAdminOrAuctioneer middleware ensures the user is either a system admin or auctioneer
func RequireAdminOrAuctioneer() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		adminRole, ok := role.(domain.AdminRole)
		if !ok || (adminRole != domain.RoleSystemAdmin && adminRole != domain.RoleAuctioneer) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetClaims retrieves JWT claims from the context
func GetClaims(c *gin.Context) (*domain.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*domain.JWTClaims)
	return jwtClaims, ok
}

// GetUserID retrieves user ID from the context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}
