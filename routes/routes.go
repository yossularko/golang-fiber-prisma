package routes

import (
	"golang-fiber-prisma/app/users"
	"golang-fiber-prisma/db"

	"github.com/gofiber/fiber/v2"
)

func HomeHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Hello World!"})
}

func setupApiRoutes(r *fiber.App, prisma *db.PrismaClient) {
	api := r.Group("/api")
	api.Get("/", HomeHandler)
	users.Routes(api.Group("/users"), prisma)
}

func SetupRoutes(r *fiber.App, prisma *db.PrismaClient) {
	setupApiRoutes(r, prisma)
}
