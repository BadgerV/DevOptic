package rbac

import (
    "context"    
    "github.com/google/uuid"
)

type Service struct {
    repo RBACRepository
}

func NewService(repo RBACRepository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateRole(ctx context.Context, name, description string) (*Role, error) {
    role := &Role{
        ID:          uuid.New(),
        Name:        name,
        Description: description,
    }
    return s.repo.CreateRole(ctx, role)
}

func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
    return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *Service) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
    return s.repo.RemoveRoleFromUser(ctx, userID, roleID)
}

func (s *Service) GetAllUsersUsernameAndID(ctx context.Context) ([]User, error) {
    return s.repo.GetAllUsersUsernameAndID(ctx)
}

func (s *Service) GetUserPermissions(ctx context.Context, userID uuid.UUID) (*UserPermissions, error) {
    roles, err := s.repo.GetUserRoles(ctx, userID)
    if err != nil {
        return nil, err
    }

    var allPermissions []Permission
    permissionMap := make(map[uuid.UUID]Permission)

    for _, role := range roles {
        permissions, err := s.repo.GetRolePermissions(ctx, role.ID)
        if err != nil {
            return nil, err
        }

        for _, permission := range permissions {
            if _, exists := permissionMap[permission.ID]; !exists {
                permissionMap[permission.ID] = permission
                allPermissions = append(allPermissions, permission)
            }
        }
    }

    return &UserPermissions{
        UserID:      userID,
        Roles:       roles,
        Permissions: allPermissions,
    }, nil
}

func (s *Service) CheckPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
    return s.repo.CheckUserPermission(ctx, userID, resource, action)
}

func (s *Service) GetAllRoles(ctx context.Context) ([]Role, error) {
    return s.repo.GetAllRoles(ctx)
}

func (s *Service) CreatePermission(ctx context.Context, resource, action, description string) (*Permission, error) {
    permission := &Permission{
        ID:          uuid.New(),
        Resource:    resource,
        Action:      action,
        Description: description,
    }
    
    return s.repo.CreatePermission(ctx, permission)
}

func (s *Service) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
    return s.repo.AssignPermissionToRole(ctx, roleID, permissionID)
}

// GetUserRoles returns all roles assigned to a user
func (s *Service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
    return s.repo.GetUserRoles(ctx, userID)
}
