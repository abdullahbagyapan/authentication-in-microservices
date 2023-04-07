package main

import (
	authGrpc "microservices/auth-service/pkg/grpc"
	"microservices/gateway-service/internal/adapters/app"
	"microservices/gateway-service/internal/adapters/redis"
	tokenGrpc "microservices/token-service/pkg/grpc"
)

func main() {

	tokenServiceClient := tokenGrpc.NewTokenServiceClient()

	authServiceClient := authGrpc.NewAuthServiceClient()

	redisClient := redis.NewAdapter()

	appAdapter := app.NewAdapter(authServiceClient, tokenServiceClient, redisClient)

	appAdapter.Run()
}
