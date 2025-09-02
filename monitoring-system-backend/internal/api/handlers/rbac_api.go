package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/badgerv/monitoring-api/internal/rbac"
)

type RBACAPI struct {
	RBAC *rbac.Service
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreatePermissionRequest struct {
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

type AssignPermissionRequest struct {
	RoleID       uuid.UUID `json:"role_id" binding:"required"`
	PermissionID uuid.UUID `json:"permission_id" binding:"required"`
}

type CheckPermissionRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Resource string    `json:"resource" binding:"required"`
	Action   string    `json:"action" binding:"required"`
}

func NewRBACHandler(r *rbac.Service) *RBACAPI {
	return &RBACAPI{RBAC: r}
}

// Health check for RBAC service
func (h *RBACAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "RBAC service OK"})
}

// --- ROLE ENDPOINTS ---

// Create a new role
func (h *RBACAPI) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	role, err := h.RBAC.CreateRole(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    role,
	})
}

// Get all roles
func (h *RBACAPI) GetAllRoles(c *gin.Context) {
	roles, err := h.RBAC.GetAllRoles(c.Request.Context())
	if err != nil {
		log.Println("get all roles error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    roles,
	})
}

// --- PERMISSION ENDPOINTS ---

// Create a new permission
func (h *RBACAPI) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	permission, err := h.RBAC.CreatePermission(c.Request.Context(), req.Resource, req.Action, req.Description)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permission created successfully",
		"data":    permission,
	})
}

// --- USER-ROLE ASSIGNMENT ENDPOINTS ---

// Assign role to user
func (h *RBACAPI) AssignRoleToUser(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	err := h.RBAC.AssignRoleToUser(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		log.Println("assign role error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role assigned successfully",
	})
}

// Remove role from user
func (h *RBACAPI) RemoveRoleFromUser(c *gin.Context) {
	var req AssignRoleRequest // Reusing the same struct since fields are identical
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	err := h.RBAC.RemoveRoleFromUser(c.Request.Context(), req.UserID, req.RoleID)
	if err != nil {
		log.Println("remove role error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role removed successfully",
	})
}

// --- ROLE-PERMISSION ASSIGNMENT ENDPOINTS ---

// Assign permission to role
func (h *RBACAPI) AssignPermissionToRole(c *gin.Context) {
	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	err := h.RBAC.AssignPermissionToRole(c.Request.Context(), req.RoleID, req.PermissionID)
	if err != nil {
		log.Println("assign permission error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission assigned successfully",
	})
}

// Fetch all users id and usernames
func (h *RBACAPI) GetAllUsersUsernameAndID(c *gin.Context) {
	users, err := h.RBAC.GetAllUsersUsernameAndID(c.Request.Context())

	if err != nil {
		log.Println("assign permission error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "users fetched successfully",
		"data" : users,
	})
}

// --- USER PERMISSION QUERIES ---

// Get user permissions (roles and permissions)
func (h *RBACAPI) GetUserPermissions(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid UUID"})
		return
	}

	permissions, err := h.RBAC.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		log.Println("get user permissions error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    permissions,
	})
}

// Check if user has specific permission
func (h *RBACAPI) CheckPermission(c *gin.Context) {
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	hasPermission, err := h.RBAC.CheckPermission(c.Request.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		log.Println("check permission error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Success",
		"has_permission": hasPermission,
	})
}

// Check permission via URL parameters (alternative endpoint)
func (h *RBACAPI) CheckUserPermission(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user UUID"})
		return
	}

	resource := c.Query("resource")
	action := c.Query("action")

	if resource == "" || action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "resource and action query parameters are required"})
		return
	}

	hasPermission, err := h.RBAC.CheckPermission(c.Request.Context(), userID, resource, action)
	if err != nil {
		log.Println("check permission error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Success",
		"user_id":        userID,
		"resource":       resource,
		"action":         action,
		"has_permission": hasPermission,
	})
}

// Check if a user is Super Admin
func (h *RBACAPI) CheckIfSuperAdmin(c *gin.Context) {
    userIDParam := c.Param("user_id")
    userID, err := uuid.Parse(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user UUID"})
        return
    }

    roles, err := h.RBAC.GetUserRoles(c.Request.Context(), userID)
    if err != nil {
        log.Println("get user roles error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    isSuperAdmin := false
    for _, role := range roles {
        if role.Name == "super admin" {
            isSuperAdmin = true
            break
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id":        userID,
        "is_super_admin": isSuperAdmin,
    })
}
