package rbac

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	RoleID     uuid.UUID `json:"role_id" db:"role_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID `json:"permission_id" db:"permission_id"`
}

// For responses
type RoleWithPermissions struct {
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

type UserPermissions struct {
	UserID      uuid.UUID    `json:"user_id"`
	Roles       []Role       `json:"roles"`
	Permissions []Permission `json:"permissions"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// Repository interface
type RBACRepository interface {
	// Role operations
	CreateRole(ctx context.Context, role *Role) (*Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	GetAllRoles(ctx context.Context) ([]Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// Permission operations
	CreatePermission(ctx context.Context, permission *Permission) (*Permission, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error)
	GetPermissionByResourceAction(ctx context.Context, resource, action string) (*Permission, error)
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	UpdatePermission(ctx context.Context, permission *Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error

	// Role-Permission operations
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)

	// User-Role operations
	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error)

	// Permission checking
	CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)
	CheckUserRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error)

	GetAllUsersUsernameAndID(ctx context.Context) ([]User, error)
}
