// Purpose: HTTP middleware — JWT auth, RBAC role enforcement, blacklist check
// Ref: ADR-002 — JWT + Redis blacklist; all protected routes must pass through AuthMiddleware
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/tigersoft/ats-go-api/config"
	pkgjwt "github.com/tigersoft/ats-go-api/pkg/jwt"
	pkgredis "github.com/tigersoft/ats-go-api/pkg/redis"
	"github.com/tigersoft/ats-go-api/pkg/response"
)

// AuthMiddleware validates the Bearer JWT and checks the Redis blacklist.
// On success, sets "user_id", "role", and "jti" in the Gin context.
func AuthMiddleware(cfg *config.Config, rdb *goredis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract Bearer token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		tokenStr := parts[1]

		// 2. Parse + validate JWT
		claims, err := pkgjwt.Parse(tokenStr, cfg.JWTSecret)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 3. Check Redis blacklist for claims.JTI
		blacklisted, err := pkgredis.IsBlacklisted(c.Request.Context(), rdb, claims.JTI)
		if err != nil || blacklisted {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 4. Set context values for downstream handlers
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("jti", claims.JTI)
		c.Set("token", tokenStr)

		c.Next()
	}
}

// RequireRole returns a middleware that enforces role-based access control.
// Pass one or more allowed roles; request is forbidden if user role is not in the list.
func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || !roleSet[roleStr] {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
