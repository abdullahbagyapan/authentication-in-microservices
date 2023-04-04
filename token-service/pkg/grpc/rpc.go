package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"microservices/token-service/internal/ports"
	pb "microservices/token-service/pkg/grpc/tokenpb"
)

const address = "localhost:9000"

func NewGRPCAdapter(app ports.AppPort) *Adapter {
	return &Adapter{api: app}
}

func NewTokenServiceClient() pb.TokenServiceClient {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error did not connect to token service %v", err)
	}

	log.Println("Successfully created token service client")
	return pb.NewTokenServiceClient(conn)

}

func (grpcA Adapter) GenerateToken(ctx context.Context, req *pb.NewTokenRequest) (*pb.NewTokenResponse, error) {
	token, err := grpcA.api.GenerateToken(req.UserId)

	if err != nil {
		return nil, status.Error(codes.Internal, "unexcepted error")
	}

	return &pb.NewTokenResponse{Token: token}, nil

}

func (grpcA Adapter) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {

	userId, isTokenValid, err := grpcA.api.ValidateToken(req.Token)

	if err != nil {
		return nil, status.Error(codes.Internal, "unexcepted error")
	}

	return &pb.ValidateTokenResponse{IsValid: isTokenValid, UserId: userId}, nil

}
