package routes

import (
	"golang-fiber-prisma/app/roles"
	"golang-fiber-prisma/app/users"
	"golang-fiber-prisma/inits"

	"github.com/gofiber/fiber/v2"
)

func HomeHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Hello World!"})
}

type TestRequest struct {
	FirstName      string `json:"first_name" validate:"min=3"`
	LastName       string `json:"last_name" validate:"required,max=10"`
	Age            uint8  `json:"age" validate:"gte=15,lte=130"`
	Email          string `json:"email" validate:"required,email"`
	Gender         string `json:"gender" validate:"oneof=male female prefer_not_to"`
	FavouriteColor string `json:"favourite_color" validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
}

func HomePostHandler(c *fiber.Ctx) error {
	var body TestRequest
	if errParse := c.BodyParser(&body); errParse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errParse.Error()})
	}

	// validation
	if err := inits.MyValidate(body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "data": body})
}

func setupApiRoutes(r *fiber.App) {
	api := r.Group("/api")
	api.Get("/", HomeHandler)
	api.Post("/", HomePostHandler)
	users.Routes(api.Group("/users"))
	roles.Routes(api.Group("/roles"))
}

func SetupRoutes(r *fiber.App) {
	setupApiRoutes(r)
}
