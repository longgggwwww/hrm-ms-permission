package grpc_clients

import (
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(connStr string) *pb.UserServiceClient {
	conn, err := grpc.NewClient(connStr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	client := pb.NewUserServiceClient(conn)
	return &client
}
