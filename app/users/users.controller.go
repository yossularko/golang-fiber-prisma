package users

import (
	"golang-fiber-prisma/lib"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func requestQueryIndex(c *fiber.Ctx) (UserQueryRequest, error) {
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
		return UserQueryRequest{}, errPageInt
	}

	perpageInt, errPerpageInt := strconv.Atoi(perPage)

	if errPerpageInt != nil {
		return UserQueryRequest{}, errPerpageInt
	}

	query := UserQueryRequest{
		Name:    q["name"],
		Email:   q["email"],
		Page:    pageInt,
		PerPage: perpageInt,
	}
	return query, nil
}

func IndexHandler(c *fiber.Ctx) error {
	q, err := requestQueryIndex(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			lib.ResponseError(lib.ResponseProps{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			}),
		)
	}

	result := GetAllUsersService(q)
	return c.Status(result.StatusCode).JSON(result)
}

func StoreHandler(c *fiber.Ctx) error {
	var user UserRequest
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	result := CreateOneService(user)
	return c.Status(result.StatusCode).JSON(result)
}
