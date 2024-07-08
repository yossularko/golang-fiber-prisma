package hasroles

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetHasRoleById(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, errInt := strconv.Atoi(id)

	if errInt != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errInt.Error()})
	}

	result := GetHasRoleByIdService(idInt)
	return c.Status(result.StatusCode).JSON(result)
}

func GetHasRolesByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")

	usrIdInt, errInt := strconv.Atoi(userId)

	if errInt != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errInt.Error()})
	}

	result := GetHasRolesByUserIdService(usrIdInt)
	return c.Status(result.StatusCode).JSON(result)
}
