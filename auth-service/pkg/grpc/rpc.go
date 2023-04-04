package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"microservices/auth-service/internal/ports"
	"microservices/auth-service/pkg/grpc/pb"
	"microservices/token-service/pkg/grpc/tokenpb"
)

const address = "localhost:9001"

func NewGRPCAdapter(app ports.AppPort, tokenService tokenpb.TokenServiceClient) *Adapter {
	return &Adapter{api: app, tokenService: tokenService}
}

func NewAuthServiceClient() pb.AuthServiceClient {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error did not connect to auth service %v", err)
	}

	log.Println("Successfully created auth service client")
	return pb.NewAuthServiceClient(conn)

}

func (grpcA Adapter) Login(ctx context.Context, req *pb.LoginRequest) (*pb.RegisterLoginResponse, error) {
	Id, isUserValid, err := grpcA.api.Login(req.Username, req.Password)

	if err != nil {
		return nil, status.Error(codes.Internal, "unexcepted error")
	}

	if !isUserValid {
		return nil, status.Error(codes.NotFound, "user is not valid")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// invoke token service
	token, err := grpcA.tokenService.GenerateToken(ctx, &tokenpb.NewTokenRequest{UserId: Id})

	if err != nil {
		log.Printf("Error requesting token service %v", err)
		return nil, status.Error(codes.Internal, "unexcepted error")

	}

	return &pb.RegisterLoginResponse{Token: token.Token}, nil

}

func (grpcA Adapter) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterLoginResponse, error) {

	userId, err := grpcA.api.Register(req.Name, req.Username, req.Password, req.Email)

	if err != nil {
		return nil, status.Error(codes.Internal, "unexcepted error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// invoke token service
	token, err := grpcA.tokenService.GenerateToken(ctx, &tokenpb.NewTokenRequest{UserId: userId})

	if err != nil {
		log.Printf("Error requesting token service %v", err)
		return nil, status.Error(codes.Internal, "unexcepted error")

	}

	return &pb.RegisterLoginResponse{Token: token.Token}, nil

}
