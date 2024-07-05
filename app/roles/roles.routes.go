package roles

import "github.com/gofiber/fiber/v2"

func Routes(r fiber.Router) {
	r.Get("/", GetAllRole)
	r.Get("/:id", GetRole)
}