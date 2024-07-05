package users

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router) {
	r.Get("/", IndexHandler)
	r.Post("/", StoreHandler)
}
