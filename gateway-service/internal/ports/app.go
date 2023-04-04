package ports

import "github.com/gofiber/fiber/v2"

type HttpPort interface {
	Run()
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	CreateComment(c *fiber.Ctx) error
	GetAllComment(c *fiber.Ctx) error
	ValidateToken(c *fiber.Ctx) error
}
