package main

import (
	authGrpc "microservices/auth-service/pkg/grpc"
	"microservices/gateway-service/internal/adapters/app"
	tokenGrpc "microservices/token-service/pkg/grpc"
)

func main() {

	tokenServiceClient := tokenGrpc.NewTokenServiceClient()

	authServiceClient := authGrpc.NewAuthServiceClient()

	appAdapter := app.NewAdapter(authServiceClient, tokenServiceClient)

	appAdapter.Run()
}
