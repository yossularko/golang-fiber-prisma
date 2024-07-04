package users

import (
	"golang-fiber-prisma/db"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, prisma *db.PrismaClient) {
	r.Get("/", func(c *fiber.Ctx) error { return IndexHandler(c, prisma) })
	r.Post("/", func(c *fiber.Ctx) error { return StoreHandler(c, prisma) })
}
