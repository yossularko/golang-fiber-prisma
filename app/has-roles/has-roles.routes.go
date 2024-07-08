package hasroles

import "github.com/gofiber/fiber/v2"

func Routes(r fiber.Router) {
	r.Get("/:id", GetHasRoleById)
	r.Get("/user/:userId", GetHasRolesByUserId)
	r.Post("/", CreateHasRole)
	r.Patch("/:id", UpdateHasRole)
	r.Delete("/:id", DeleteHasRole)
}
