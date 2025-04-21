package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/proto/entpb"
	"github.com/longgggwww/hrm-ms-permission/internal/handlers"
	"google.golang.org/grpc"
)

func startGRPCServer(client *ent.Client) {
	perm := entpb.NewPermService(client)
	permGroup := entpb.NewPermGroupService(client)

	server := grpc.NewServer()
	fmt.Println("Starting gRPC server on port 5000...")

	entpb.RegisterPermServiceServer(server, perm)
	entpb.RegisterPermGroupServiceServer(server, permGroup)

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

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	client, err := ent.Open("postgres", "host=localhost port=5433 user=root dbname=permission password=123456 sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	go startHTTPServer(client)
	startGRPCServer(client)
}
