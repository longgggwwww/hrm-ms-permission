package main

import (
	"context"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/proto/entpb"
	"google.golang.org/grpc"
)

func main() {
	// Initialize an ent client.
	client, err := ent.Open("postgres", "host=localhost port=5432 user=root dbname=permission password=123456 sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Run the migration tool (creating tables, etc).
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// permGroup := client.PermGroup.Create().
	// 	SetCode("PERM_GROUP_1").
	// 	SetName("Permission Group 1").SaveX(ctx)

	// fmt.Println(permGroup)

	// perm := client.Perm.Create().
	// 	SetCode("PERM_1").
	// 	SetName("Permission 1").
	// 	SetDescription("This is permission 1").SetGroupID(permGroup.ID).SaveX(ctx)

	// fmt.Println(perm)

	// Initialize the generated User service.
	svc := entpb.NewPermService(client)
	svcGroup := entpb.NewPermGroupService(client)

	// Create a new gRPC server (you can wire multiple services to a single server).
	server := grpc.NewServer()

	fmt.Println("Starting gRPC server on port 5000...")

	// Register the User service with the server.
	entpb.RegisterPermServiceServer(server, svc)
	entpb.RegisterPermGroupServiceServer(server, svcGroup)

	// Open port 5000 for listening to traffic.
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	}

	// Listen for traffic indefinitely.
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}
