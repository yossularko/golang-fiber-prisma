package users

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router) {
	r.Get("/", GetAllUser)
	r.Get("/:id", GetUser)
	r.Post("/", CreateUser)
}
