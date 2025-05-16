package entpb

import (
	"context"

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
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	_, err := s.client.UserPerm.Delete().Where(userperm.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user perms: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteUserRolesByUserID deletes all UserRoles by user_id.
func (s *ExtService) DeleteUserRolesByUserID(ctx context.Context, req *DeleteUserRolesByUserIDRequest) (*emptypb.Empty, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	_, err := s.client.UserRole.Delete().Where(userrole.UserID(req.GetUserId())).Exec(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user roles: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// NewExtService returns a new ExtService.
func NewExtService(client *ent.Client) *ExtService {
	return &ExtService{
		client: client,
	}
}
