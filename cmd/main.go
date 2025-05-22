package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/proto/entpb"
	"github.com/longgggwwww/hrm-ms-permission/internal/grpc_clients"
	"github.com/longgggwwww/hrm-ms-permission/internal/handlers"
	"google.golang.org/grpc"
)

func registerGRPCServices(srv *grpc.Server, cli *ent.Client) {
	entpb.RegisterPermServiceServer(srv, entpb.NewPermService(cli))
	entpb.RegisterPermGroupServiceServer(srv, entpb.NewPermGroupService(cli))
	entpb.RegisterRoleServiceServer(srv, entpb.NewRoleService(cli))
	entpb.RegisterUserRoleServiceServer(srv, entpb.NewUserRoleService(cli))
	entpb.RegisterUserPermServiceServer(srv, entpb.NewUserPermService(cli))
	entpb.RegisterExtServiceServer(srv, entpb.NewExtService(cli))
}

func startGRPCServer(cli *ent.Client) {
	srv := grpc.NewServer()
	registerGRPCServices(srv, cli)

	log.Println("gRPC server started on :5000")
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed listening: %v", err)
	}

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("gRPC server ended: %v", err)
	}
}

func startHTTPServer(cli *ent.Client) {
	r := gin.Default()

	userSrv := grpc_clients.NewUserClient(os.Getenv("USER_SERVICE_URL"))

	// Đăng ký các route cho HTTP server
	handlersList := []struct {
		register func(*gin.Engine)
	}{
		{handlers.NewPermGroupHandler(cli).RegisterRoutes},
		{handlers.NewPermHandler(cli, userSrv).RegisterRoutes},
		{handlers.NewRoleHandler(cli, userSrv).RegisterRoutes},
	}
	for _, h := range handlersList {
		h.register(r)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	cli, err := ent.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer cli.Close()

	go startHTTPServer(cli)
	startGRPCServer(cli)
}
