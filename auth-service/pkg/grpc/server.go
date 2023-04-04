package grpc

import (
	"google.golang.org/grpc"
	"log"
	"microservices/auth-service/internal/ports"
	"microservices/auth-service/pkg/grpc/pb"
	"microservices/token-service/pkg/grpc/tokenpb"
	"net"
)

type Adapter struct {
	api ports.AppPort
	pb.UnimplementedAuthServiceServer
	tokenService tokenpb.TokenServiceClient
}

func (grpcA Adapter) Run() {

	listen, err := net.Listen("tcp", ":9001")

	if err != nil {
		log.Fatalf("failed to listening %v", err)
	}

	authService := grpcA

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Println("Auth server is serving at : ", listen.Addr())

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serving grpc server %v", err)
	}

}
