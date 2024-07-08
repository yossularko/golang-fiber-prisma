package hasroles

import (
	"context"
	"golang-fiber-prisma/app/users"
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/inits"
	"golang-fiber-prisma/lib"

	"github.com/gofiber/fiber/v2"
)

func GetHasRoleByIdService(id int) lib.ResponseData {
	data, err := inits.Prisma.HasRoles.
		FindUnique(db.HasRoles.ID.Equals(id)).
		Select(
			db.HasRoles.ID.Field(),
			db.HasRoles.CreatedAt.Field(),
			db.HasRoles.UpdatedAt.Field(),
		).
		With(
			db.HasRoles.Role.Fetch().Select(
				db.Roles.ID.Field(),
				db.Roles.Name.Field(),
			),
			db.HasRoles.User.Fetch().Select(
				db.User.ID.Field(),
				db.User.Name.Field(),
				db.User.Email.Field(),
			),
		).
		Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: data})
}

func GetHasRolesByUserIdService(userId int) lib.ResponseData {
	user, errUsr := users.GetUserOne(userId, "")

	if errUsr != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errUsr.Error()})
	}

	data, err := inits.Prisma.HasRoles.
		FindMany(
			db.HasRoles.UserID.Equals(user.ID),
		).
		Select(
			db.HasRoles.ID.Field(),
			db.HasRoles.CreatedAt.Field(),
			db.HasRoles.UpdatedAt.Field(),
		).
		With(
			db.HasRoles.Role.Fetch().Select(
				db.Roles.ID.Field(),
				db.Roles.Name.Field(),
			),
		).
		Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	response := []HasRolesWithRole{}

	for _, hasRole := range data {
		response = append(response, HasRolesWithRole{
			ID:        hasRole.ID,
			CreatedAt: hasRole.CreatedAt,
			UpdatedAt: hasRole.UpdatedAt,
			Roles: RolesInResponse{
				ID:   hasRole.Role().ID,
				Name: hasRole.Role().Name,
			},
		})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}
