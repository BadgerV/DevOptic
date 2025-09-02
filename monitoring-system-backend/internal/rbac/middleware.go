package rbac

import (
	"fmt"
	"net/http"

	"github.com/badgerv/monitoring-api/internal/auth"
	"github.com/badgerv/monitoring-api/internal/storage"
	"github.com/gin-gonic/gin"
)

const AuthContextKey = "auth_context"

func (s *Service) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := auth.GetAuthContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		hasPermission, err := s.CheckPermission(c.Request.Context(), authContext.User.ID, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Service) RequireRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := auth.GetAuthContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userPerms, err := s.GetUserPermissions(c.Request.Context(), authContext.User.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role check failed"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range userPerms.Roles {
			for _, requiredRole := range roleNames {
				if role.Name == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func (s *Service) RequireRoleForWebsocket(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := storage.NewDB()

		userRepo := auth.NewPostgresUserRepository(db)
		sessionRepo := auth.NewPostgresSessionRepository(db)

		authService := auth.NewService(userRepo, sessionRepo)
		// Extract token from query param instead of Authorization header
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		fmt.Println("Getting Here, gotten the token from query", token)

		// Validate the token (replace with your own validation logic)
		authContext, err := authService.ValidateToken(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		fmt.Println("Validated the token, auth context", authContext)

		// Fetch user permissions
		userPerms, err := s.GetUserPermissions(c.Request.Context(), authContext.User.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role check failed"})
			c.Abort()
			return
		}

		fmt.Println("User permissions - ", userPerms)


		// Check roles
		hasRole := false
		for _, role := range userPerms.Roles {
			for _, requiredRole := range roleNames {
				if role.Name == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		fmt.Println("These are the roles the user has - ", hasRole)

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
			c.Abort()
			return
		}

		// Save authContext in Gin context for use later in WebSocket handler
		c.Set(AuthContextKey, authContext)

		c.Next()
	}
}
