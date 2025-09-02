package rbac

import (
	"context"
	"fmt"
	"time"

	"github.com/badgerv/monitoring-api/internal/storage" // Update this import path to match your project
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db *storage.DB
}

func NewPostgresRepository(db *storage.DB) RBACRepository {
	return &PostgresRepository{db: db}
}

// Role operations
func (r *PostgresRepository) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	query := `
        INSERT INTO roles (id, name, description, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, description, created_at`

	var created Role
	err := r.db.Pool.QueryRow(ctx, query,
		role.ID, role.Name, role.Description, time.Now()).
		Scan(&created.ID, &created.Name, &created.Description, &created.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &created, nil
}

func (r *PostgresRepository) GetAllUsersUsernameAndID(ctx context.Context) ([]User, error) {
	query := `
        SELECT id, username
        FROM users
        ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return users, nil
}

func (r *PostgresRepository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`

	var role Role
	err := r.db.Pool.QueryRow(ctx, query, name).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

func (r *PostgresRepository) GetAllRoles(ctx context.Context) ([]Role, error) {
	query := `SELECT id, name, description, created_at FROM roles ORDER BY name`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating roles: %w", rows.Err())
	}

	return roles, nil
}

// Permission operations
func (r *PostgresRepository) CreatePermission(ctx context.Context, permission *Permission) (*Permission, error) {
	query := `
        INSERT INTO permissions (id, resource, action, description, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, resource, action, description, created_at`

	var created Permission
	err := r.db.Pool.QueryRow(ctx, query,
		permission.ID, permission.Resource, permission.Action,
		permission.Description, time.Now()).
		Scan(&created.ID, &created.Resource, &created.Action,
			&created.Description, &created.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &created, nil
}

// User-Role operations
func (r *PostgresRepository) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
        INSERT INTO user_roles (user_id, role_id, assigned_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, role_id) DO NOTHING
        RETURNING user_id`

	var insertedUserID uuid.UUID
	err := r.db.Pool.QueryRow(ctx, query, userID, roleID).Scan(&insertedUserID)
	if err != nil {
		// if no row was inserted, Scan will fail with pgx.ErrNoRows
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("user already has the role")
		}
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	query := `
        SELECT r.id, r.name, r.description, r.created_at FROM roles r
        JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1
        ORDER BY r.name`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating user roles: %w", rows.Err())
	}

	return roles, nil
}

func (r *PostgresRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error) {
	query := `
        SELECT p.id, p.resource, p.action, p.description, p.created_at FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        WHERE rp.role_id = $1
        ORDER BY p.resource, p.action`

	rows, err := r.db.Pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		err := rows.Scan(&permission.ID, &permission.Resource, &permission.Action,
			&permission.Description, &permission.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating role permissions: %w", rows.Err())
	}

	return permissions, nil
}

func (r *PostgresRepository) CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	query := `
        SELECT COUNT(*) > 0 FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        JOIN user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = $1 AND p.resource = $2 AND p.action = $3`

	var hasPermission bool
	err := r.db.Pool.QueryRow(ctx, query, userID, resource, action).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}

	return hasPermission, nil
}

func (r *PostgresRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*Role, error) {
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`

	var role Role
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

func (r *PostgresRepository) UpdateRole(ctx context.Context, role *Role) error {
	query := `UPDATE roles SET name = $2, description = $3 WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query, role.ID, role.Name, role.Description)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	// Use transaction for cascading deletes
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Remove role permissions
	_, err = tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove role permissions: %w", err)
	}

	// Remove user role assignments
	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE role_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove user role assignments: %w", err)
	}

	// Delete the role
	_, err = tx.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (*Permission, error) {
	query := `SELECT id, resource, action, description, created_at FROM permissions WHERE id = $1`

	var permission Permission
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&permission.ID, &permission.Resource, &permission.Action,
			&permission.Description, &permission.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

func (r *PostgresRepository) GetPermissionByResourceAction(ctx context.Context, resource, action string) (*Permission, error) {
	query := `SELECT id, resource, action, description, created_at FROM permissions WHERE resource = $1 AND action = $2`

	var permission Permission
	err := r.db.Pool.QueryRow(ctx, query, resource, action).
		Scan(&permission.ID, &permission.Resource, &permission.Action,
			&permission.Description, &permission.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

func (r *PostgresRepository) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	query := `SELECT id, resource, action, description, created_at FROM permissions ORDER BY resource, action`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		err := rows.Scan(&permission.ID, &permission.Resource, &permission.Action,
			&permission.Description, &permission.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", rows.Err())
	}

	return permissions, nil
}

func (r *PostgresRepository) UpdatePermission(ctx context.Context, permission *Permission) error {
	query := `UPDATE permissions SET resource = $2, action = $3, description = $4 WHERE id = $1`

	_, err := r.db.Pool.Exec(ctx, query, permission.ID, permission.Resource, permission.Action, permission.Description)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	// Use transaction for cascading deletes
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Remove role permission assignments
	_, err = tx.Exec(ctx, `DELETE FROM role_permissions WHERE permission_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove role permission assignments: %w", err)
	}

	// Delete the permission
	_, err = tx.Exec(ctx, `DELETE FROM permissions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	_, err := r.db.Pool.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	// Step 1: Validate the role assignment exists
	var exists bool
	checkQuery := `
        SELECT EXISTS (
            SELECT 1 FROM user_roles
            WHERE user_id = $1 AND role_id = $2
        )
    `
	if err := r.db.Pool.QueryRow(ctx, checkQuery, userID, roleID).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check role assignment: %w", err)
	}

	if !exists {
		return fmt.Errorf("cannot remove role: user %s does not have role %s", userID, roleID)
	}

	// Step 2: Proceed with deletion
	delQuery := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	cmdTag, err := r.db.Pool.Exec(ctx, delQuery, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	// Safety check: should always affect exactly 1 row
	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("unexpected result: no rows deleted while removing role %s from user %s", roleID, userID)
	}

	return nil
}

func (r *PostgresRepository) GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM user_roles WHERE role_id = $1`

	rows, err := r.db.Pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query role users: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		err := rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating role users: %w", rows.Err())
	}

	return userIDs, nil
}

func (r *PostgresRepository) CheckUserRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	query := `
        SELECT COUNT(*) > 0 FROM roles r
        JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1 AND r.name = $2`

	var hasRole bool
	err := r.db.Pool.QueryRow(ctx, query, userID, roleName).Scan(&hasRole)
	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return hasRole, nil
}

func (r *PostgresRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
        ON CONFLICT (role_id, permission_id) DO NOTHING`

	_, err := r.db.Pool.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return nil
}
