package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/proto/entpb"
	"github.com/longgggwww/hrm-ms-permission/internal/handlers"
	"google.golang.org/grpc"
)

func startGRPCServer(client *ent.Client) {
	perm := entpb.NewPermService(client)
	permGroup := entpb.NewPermGroupService(client)
	role := entpb.NewRoleService(client)

	server := grpc.NewServer()
	fmt.Println("Starting gRPC server on port 5000...")

	entpb.RegisterPermServiceServer(server, perm)
	entpb.RegisterPermGroupServiceServer(server, permGroup)
	entpb.RegisterRoleServiceServer(server, role)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}

func startHTTPServer(client *ent.Client) {
	r := gin.Default()

	r.GET("/perms", handlers.GetPermsHandler(client))
	r.GET("/roles", handlers.GetRolesHandler(client))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DB_URI := os.Getenv("DB_URI")
	if DB_URI == "" {
		log.Fatalf("DB_URI is not set in the environment variables")
	}

	client, err := ent.Open("postgres", DB_URI)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	go startHTTPServer(client)
	startGRPCServer(client)
}
