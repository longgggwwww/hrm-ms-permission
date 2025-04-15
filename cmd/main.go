package main

import (
	"context"
	"log"
	"net"

	pb "github.com/longgggwww/hrm-ms-permission/api"

	"google.golang.org/grpc"
)

type service struct {
	pb.UnimplementedPermissionServer
}

// Implement the CheckPermission method
func (s *service) CheckPermission(ctx context.Context, req *pb.PermissionRequest) (*pb.PermissionResponse, error) {
	// Example logic: Allow all requests for now
	return &pb.PermissionResponse{Allowed: true}, nil
}

func main() {
	// Set up a listener on a specific port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register your service implementation with the gRPC server
	pb.RegisterPermissionServer(grpcServer, &service{})

	log.Println("gRPC server is running on port 50051...")

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
