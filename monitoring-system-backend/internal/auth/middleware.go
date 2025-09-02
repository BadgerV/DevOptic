package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const AuthContextKey = "auth_context"

func (s *Service) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        authContext, err := s.ValidateToken(c.Request.Context(), token)
        fmt.Println(err)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        fmt.Println("Token Validated", token)


        c.Set(AuthContextKey, authContext)
        c.Set("userID", authContext.User.ID)

        c.Next()
    }
}

func GetAuthContext(c *gin.Context) (*AuthContext, bool) {
    authContext, exists := c.Get(AuthContextKey)
    if !exists {
        return nil, false
    }
    
    ctx, ok := authContext.(*AuthContext)
    return ctx, ok
}