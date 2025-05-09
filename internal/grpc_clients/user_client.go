package grpc_clients

import (
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient() (*pb.UserServiceClient, error) {
	conn, err := grpc.NewClient("user_app:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &client, nil
}
