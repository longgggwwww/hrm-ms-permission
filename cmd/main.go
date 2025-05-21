package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/proto/entpb"
	"github.com/longgggwwww/hrm-ms-permission/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-permission/internal/handlers"
	"google.golang.org/grpc"
)

func startGRPCServer(cli *ent.Client) {
	perm := entpb.NewPermService(cli)
	permGroup := entpb.NewPermGroupService(cli)
	role := entpb.NewRoleService(cli)
	userRole := entpb.NewUserRoleService(cli)
	userPerm := entpb.NewUserPermService(cli)

	server := grpc.NewServer()
	fmt.Println("Starting gRPC server on port 5000...")

	entpb.RegisterPermServiceServer(server, perm)
	entpb.RegisterPermGroupServiceServer(server, permGroup)
	entpb.RegisterRoleServiceServer(server, role)
	entpb.RegisterUserRoleServiceServer(server, userRole)
	entpb.RegisterUserPermServiceServer(server, userPerm)
	entpb.RegisterPermExtServiceServer(server, entpb.NewPermExtService(cli))

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}

func startHTTPServer(cli *ent.Client, roleHandler *handlers.RoleHandler) {
	r := gin.Default()

	permGroup := handlers.PermGroupHandler{Client: cli}
	permGroup.RegisterRoutes(r)

	perm := handlers.PermHandler{Client: cli}
	perm.RegisterRoutes(r)

	roleHandler.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	// Initialize gRPC clients
	userClient, err := grpc_clients.NewUserClient()
	if err != nil {
		log.Fatalf("failed to initialize user client: %v", err)
	}

	roleHandler := handlers.NewRoleHandler(cli, userClient)

	// Pass roleHandler to HTTP server
	go startHTTPServer(cli, roleHandler)
	startGRPCServer(cli)
}
