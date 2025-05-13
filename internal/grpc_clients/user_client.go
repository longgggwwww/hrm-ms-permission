package grpc_clients

import (
	"log"
	"os"

	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient() (*pb.UserServiceClient, error) {
	url := os.Getenv("USER_SERVICE")
	if url == "" {
		return nil, log.Output(2, "USER_SERVICE environment variable is not set")
	}
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &client, nil
}
