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

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}
