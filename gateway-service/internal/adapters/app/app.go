package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"microservices/auth-service/pkg/grpc/pb"
	"microservices/gateway-service/internal/ports"
	"microservices/token-service/pkg/grpc/tokenpb"
)

type Adapter struct {
	authService  pb.AuthServiceClient
	tokenService tokenpb.TokenServiceClient
	redis        ports.RedisPort
}

type registerRequest struct {
	name     string
	username string
	password string
	email    string
}

type loginRequest struct {
	username string
	password string
}

func NewAdapter(authService pb.AuthServiceClient, tokenService tokenpb.TokenServiceClient, redis ports.RedisPort) *Adapter {
	return &Adapter{
		authService:  authService,
		tokenService: tokenService,
		redis:        redis,
	}
}

func (A Adapter) Run() {
	app := fiber.New()

	// whitelist
	app.Use("/login", A.Login)
	app.Use("/register", A.Register)

	app.Use(A.ValidateToken)

	app.Get("/home", func(ctx *fiber.Ctx) error {
		return ctx.JSON("hello world")
	})

	log.Fatal(app.Listen(":8080"))
}

func (A Adapter) Login(c *fiber.Ctx) error {

	var req = loginRequest{}

	err := c.BodyParser(&req)

	if err != nil {
		c.Status(500)

		return c.JSON("parse error")
	}

	grpcCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := A.authService.Login(grpcCtx, &pb.LoginRequest{
		Username: req.username,
		Password: req.password,
	})

	if err != nil {
		c.Status(500)

		return c.JSON("Login error")
	}

	// cache token
	go func() {
		err := A.redis.SetToken(resp.Token)
		if err != nil {
			log.Printf("Error setting token to Redis %v", err)
		}

	}()

	c.Status(200)
	rtnMap := map[string]string{}

	rtnMap["token"] = resp.Token

	return c.JSON(rtnMap)

}

func (A Adapter) Register(c *fiber.Ctx) error {
	var req = registerRequest{}

	err := c.BodyParser(&req)

	if err != nil {
		c.Status(500)

		return c.JSON("parse error")
	}

	grpcCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := A.authService.Register(grpcCtx, &pb.RegisterRequest{
		Name:     req.name,
		Username: req.username,
		Email:    req.email,
		Password: req.password,
	})

	if err != nil {
		c.Status(500)

		return c.JSON("Register error")
	}
	// cache token
	go func() {
		err := A.redis.SetToken(resp.Token)
		if err != nil {
			log.Printf("Error setting token to Redis %v", err)
		}
	}()

	token := resp.Token

	c.Status(200)
	rtnMap := map[string]string{}

	rtnMap["token"] = token

	return c.JSON(rtnMap)
}

func (A Adapter) ValidateToken(ctx *fiber.Ctx) error {

	authValue := string(ctx.Request().Header.Peek("Authorization"))

	if authValue == "" {
		ctx.Status(403)

		return ctx.JSON("Undefined token")
	}

	token := parseToken(authValue)

	// first look cached tokens
	err := A.redis.GetToken(token)

	if err == nil {
		return ctx.Next()
	} else {
		log.Printf("Error caching token %v", err)
	}

	grpcCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := A.tokenService.ValidateToken(grpcCtx, &tokenpb.ValidateTokenRequest{Token: token})

	if err != nil {
		ctx.Status(403)

		return ctx.JSON("Invalid token")
	}

	if !resp.IsValid {
		ctx.Status(403)

		return ctx.JSON("Token is not valid")
	}
	return ctx.Next()
}

func parseToken(value string) string {

	token := value[7:]

	return token

}
