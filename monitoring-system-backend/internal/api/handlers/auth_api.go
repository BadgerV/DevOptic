package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/badgerv/monitoring-api/internal/auth"
)

type AuthAPI struct {
	Auth              *auth.Service
	AuthLocalProvider auth.AuthProvider
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest defines the request payload for changing a user's password.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// DeliveryEmailRequest defines the request payload for setting a user's delivery email.
type DeliveryEmailRequest struct {
	DeliveryEmail string `json:"delivery_email"`
}

func NewAuthHandler(a *auth.Service) *AuthAPI {
	return &AuthAPI{Auth: a}
}

// Health check for auth service
func (h *AuthAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Auth service OK"})
}

// --- AUTH ENDPOINTS ---

// Register a new user
func (h *AuthAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	createdUser, err := h.Auth.Register(c.Request.Context(), "local", (*auth.RegisterRequest)(&req))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    createdUser,
	})
}

func (h *AuthAPI) Login(c *gin.Context) {
	var req *LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	loginResp, err := h.Auth.Login(c.Request.Context(), "local", (*auth.LoginRequest)(req))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   loginResp.Token,
		"user":    loginResp.User,
	})
}

// Logout user (invalidates session)
func (h *AuthAPI) Logout(c *gin.Context) {
	// token comes from cookie or header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing Authorization header"})
		return
	}

	err := h.Auth.Logout(c.Request.Context(), token)
	if err != nil {
		log.Println("logout error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// Validate session token
func (h *AuthAPI) Validate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
		c.Abort()
		return
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing Authorization header"})
		return
	}

	user, err := h.Auth.ValidateToken(c.Request.Context(), token)
	if err != nil {
		fmt.Println("validate error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token valid",
		"user":    user,
	})
}

// Get user by ID
func (h *AuthAPI) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid UUID"})
		return
	}

	user, err := h.Auth.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		log.Println("get user error:", err)
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    user,
	})
}

// Handler function for password change (unchanged, as confirmed correct)
func (h *AuthAPI) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	fmt.Println(req)

	// Assuming auth context or middleware provides userID
	userID := c.MustGet("userID").(uuid.UUID) // Adjust based on auth middleware

	err := h.Auth.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

// SetDeliveryEmail handles the HTTP request to set a user's delivery email.
func (h *AuthAPI) SetDeliveryEmail(c *gin.Context) {
	var req DeliveryEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	// --- Validation ---
	email := strings.TrimSpace(req.DeliveryEmail)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "delivery email is required"})
		return
	}
	if len(email) > 254 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is too long"})
		return
	}
	if !isValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email format"})
		return
	}

	// Assuming auth middleware provides userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not authenticated"})
		return
	}

	typedUserID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user ID"})
		return
	}

	// Save delivery email
	err := h.Auth.SetDeliveryEmail(c.Request.Context(), typedUserID, email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// For production, you might want to use a library like go-playground/validator.
func isValidEmail(email string) bool {
	// Simplified RFC5322 regex (enough for most cases)
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}


// GetDeliveryEmail handles the HTTP request to fetch a user's delivery email.
func (h *AuthAPI) GetDeliveryEmail(c *gin.Context) {
	// Assuming auth middleware provides userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user not authenticated"})
		return
	}

	typedUserID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid user ID"})
		return
	}

	deliveryEmail, err := h.Auth.GetDeliveryEmail(c.Request.Context(), typedUserID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Success",
		"delivery_email": deliveryEmail,
	})

}
