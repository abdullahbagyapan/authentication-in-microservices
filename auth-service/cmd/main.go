package main

import (
	"microservices/auth-service/internal/adapters/app"
	"microservices/auth-service/internal/adapters/core"
	"microservices/auth-service/internal/adapters/db"
	"microservices/auth-service/internal/adapters/rabbitmq"
	"microservices/auth-service/pkg/grpc"
	tokenGrpc "microservices/token-service/pkg/grpc"
)

func main() {

	dbAdapter := db.NewAdapter()

	msgBroker := rabbitmq.NewMsgBrokerAdapter()

	coreAdapter := core.NewAdapter()

	appAdapter := app.NewAdapter(coreAdapter, msgBroker, dbAdapter)

	tokenServiceClient := tokenGrpc.NewTokenServiceClient()

	grpcAdapter := grpc.NewGRPCAdapter(appAdapter, tokenServiceClient)

	grpcAdapter.Run()

}
