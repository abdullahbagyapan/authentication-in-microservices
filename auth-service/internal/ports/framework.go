package ports

import (
	"context"
	"microservices/auth-service/pkg/grpc/pb"
)

type User struct {
	Id       string
	Name     string
	Username string
	Password string
	Email    string
}

type MsgBrokerUserInfo struct {
	Name  string
	Email string
}

type MsgBrokerPort interface {
	PublishMessage(user *MsgBrokerUserInfo) error
}

type DbPort interface {
	FindByUsername(username string) (*User, error)
	SaveUser(user *User) error
	CloseDbConnection()
}

type GRPC interface {
	Run()
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.RegisterLoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterLoginResponse, error)
}
