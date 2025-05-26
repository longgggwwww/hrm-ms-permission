package entpb

import (
	"context"

	"github.com/google/uuid"
	ent "github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/userperm"
	"github.com/longgggwwww/hrm-ms-permission/ent/userrole"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ExtService implements ExtServiceServer.
type ExtService struct {
	client *ent.Client
	UnimplementedExtServiceServer
}

// DeleteUserPermsByUserID deletes all UserPerms by user_id.
func (s *ExtService) DeleteUserPermsByUserID(ctx context.Context, req *DeleteUserPermsByUserIDRequest) (*emptypb.Empty, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 DeleteUserPermsByUserID: user_id is required")
	}
	_, err := s.client.UserPerm.Delete().Where(userperm.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 DeleteUserPermsByUserID: failed to delete user perms: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteUserRolesByUserID deletes all UserRoles by user_id.
func (s *ExtService) DeleteUserRolesByUserID(ctx context.Context, req *DeleteUserRolesByUserIDRequest) (*emptypb.Empty, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 DeleteUserRolesByUserID: user_id is required")
	}
	_, err := s.client.UserRole.Delete().Where(userrole.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 DeleteUserRolesByUserID: failed to delete user roles: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// UpdateUserPerms updates user permissions by deleting old ones and creating new ones.
func (s *ExtService) UpdateUserPerms(ctx context.Context, req *UpdateUserPermsRequest) (*UpdateUserPermsResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 UpdateUserPerms: user_id is required")
	}
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 UpdateUserPerms: failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Delete old permissions
	_, err = tx.UserPerm.Delete().Where(userperm.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#3 UpdateUserPerms: failed to delete user perms: %v", err)
	}
	// Add new permissions
	for _, permIDStr := range req.GetPermIds() {
		permID, perr := uuid.Parse(permIDStr)
		if perr != nil {
			return nil, status.Errorf(codes.InvalidArgument, "#4 UpdateUserPerms: invalid perm_id: %v", perr)
		}
		_, err = tx.UserPerm.Create().SetUserID(req.GetUserId()).SetPermID(permID).Save(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "#5 UpdateUserPerms: failed to add user perm: %v", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "#6 UpdateUserPerms: failed to commit transaction: %v", err)
	}
	return &UpdateUserPermsResponse{Success: true}, nil
}

// UpdateUserRoles updates user roles by deleting old ones and creating new ones.
func (s *ExtService) UpdateUserRoles(ctx context.Context, req *UpdateUserRolesRequest) (*UpdateUserRolesResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 UpdateUserRoles: user_id is required")
	}
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 UpdateUserRoles: failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Delete old roles
	_, err = tx.UserRole.Delete().Where(userrole.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#3 UpdateUserRoles: failed to delete user roles: %v", err)
	}
	// Add new roles
	for _, roleIDStr := range req.GetRoleIds() {
		roleID, rerr := uuid.Parse(roleIDStr)
		if rerr != nil {
			return nil, status.Errorf(codes.InvalidArgument, "#4 UpdateUserRoles: invalid role_id: %v", rerr)
		}
		_, err = tx.UserRole.Create().SetUserID(req.GetUserId()).SetRoleID(roleID).Save(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "#5 UpdateUserRoles: failed to add user role: %v", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "#6 UpdateUserRoles: failed to commit transaction: %v", err)
	}
	return &UpdateUserRolesResponse{Success: true}, nil
}

// GetUserPerms returns all permissions of a user.
func (s *ExtService) GetUserPerms(ctx context.Context, req *GetUserPermsRequest) (*GetUserPermsResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 GetUserPerms: user_id is required")
	}
	userPerms, err := s.client.UserPerm.Query().Where(userperm.UserID(req.GetUserId())).WithPerm().All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 GetUserPerms: failed to query user perms: %v", err)
	}
	perms := make([]*PermExt, 0, len(userPerms))
	for _, up := range userPerms {
		perm := up.Edges.Perm
		if perm == nil {
			continue
		}
		perms = append(perms, &PermExt{
			Id:          perm.ID[:],
			Code:        perm.Code,
			Name:        perm.Name,
			Description: stringToStringValue(perm.Description),
		})
	}
	return &GetUserPermsResponse{Perms: perms}, nil
}

// GetUserRoles returns all roles of a user.
func (s *ExtService) GetUserRoles(ctx context.Context, req *GetUserRolesRequest) (*GetUserRolesResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "#1 GetUserRoles: user_id is required")
	}
	userRoles, err := s.client.UserRole.Query().Where(userrole.UserID(req.GetUserId())).WithRole().All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "#2 GetUserRoles: failed to query user roles: %v", err)
	}
	roles := make([]*RoleExt, 0, len(userRoles))
	for _, ur := range userRoles {
		role := ur.Edges.Role
		if role == nil {
			continue
		}
		var createdAt, updatedAt *timestamppb.Timestamp
		if role.CreatedAt != nil {
			createdAt = timestamppb.New(*role.CreatedAt)
		}
		if role.UpdatedAt != nil {
			updatedAt = timestamppb.New(*role.UpdatedAt)
		}
		roles = append(roles, &RoleExt{
			Id:          role.ID[:],
			Code:        role.Code,
			Name:        role.Name,
			Color:       stringToStringValue(role.Color),
			Description: stringToStringValue(role.Description),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}
	return &GetUserRolesResponse{Roles: roles}, nil
}

// Helper for google.protobuf.StringValue
func stringToStringValue(s string) *wrapperspb.StringValue {
	if s == "" {
		return nil
	}
	return wrapperspb.String(s)
}

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}
