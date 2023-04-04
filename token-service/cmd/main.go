package main

import (
	"microservices/token-service/internal/adapters/app"
	"microservices/token-service/internal/adapters/core"
	"microservices/token-service/pkg/grpc"
)

func main() {

	coreAdapter := core.NewCoreAdapter()

	appAdapter := app.NewAppAdapter(coreAdapter)

	grpcAdapter := grpc.NewGRPCAdapter(appAdapter)

	grpcAdapter.Run()

}
