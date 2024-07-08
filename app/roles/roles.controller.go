package roles

import (
	"golang-fiber-prisma/lib"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func requestQueryGetAll(c *fiber.Ctx) (RoleQueryRequest, error) {
	q := c.Queries()
	perPage := q["per_page"]
	page := q["page"]
	if page == "" {
		page = "1"
	}
	if perPage == "" {
		perPage = "10"
	}

	pageInt, errPageInt := strconv.Atoi(page)

	if errPageInt != nil {
		return RoleQueryRequest{}, errPageInt
	}

	perpageInt, errPerpageInt := strconv.Atoi(perPage)

	if errPerpageInt != nil {
		return RoleQueryRequest{}, errPerpageInt
	}

	query := RoleQueryRequest{
		Name:    q["name"],
		Page:    pageInt,
		PerPage: perpageInt,
	}
	return query, nil
}

func GetAllRole(c *fiber.Ctx) error {
	q, err := requestQueryGetAll(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			lib.ResponseError(lib.ResponseProps{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			}),
		)
	}

	result := GetAllRolesService(q)
	return c.Status(result.StatusCode).JSON(result)
}

func GetRole(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, errInt := strconv.Atoi(id)

	if errInt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errInt.Error()})
	}

	result := GetRoleByIdService(idInt)
	return c.Status(result.StatusCode).JSON(result)
}

func CreateRole(c *fiber.Ctx) error {
	var body RoleRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result := CreateOneService(body)
	return c.Status(result.StatusCode).JSON(result)
}

func UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, errInt := strconv.Atoi(id)

	if errInt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errInt.Error()})
	}

	var body RoleRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result := UpdateOneService(idInt, body)
	return c.Status(result.StatusCode).JSON(result)
}

func DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	idInt, errInt := strconv.Atoi(id)

	if errInt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errInt.Error()})
	}

	result := DeleteOneService(idInt)
	return c.Status(result.StatusCode).JSON(result)
}
