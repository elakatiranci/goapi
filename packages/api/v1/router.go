package v1

import (
	"github.com/gofiber/fiber/v2"
	e "zeelso.com/api/v1/endpoints"
)

func New(a *fiber.App) {
	// Register routers
	e.RegisterUserRouter(a)
}
