package hasroles

import (
	"context"
	"golang-fiber-prisma/app/roles"
	"golang-fiber-prisma/app/users"
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/inits"
	"golang-fiber-prisma/lib"

	"github.com/gofiber/fiber/v2"
)

func findByUserAndRole(userId int, roleId int) (*db.HasRolesModel, error) {
	data, err := inits.Prisma.HasRoles.
		FindFirst(
			db.HasRoles.UserID.Equals(userId),
			db.HasRoles.RoleID.Equals(roleId),
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
		return &db.HasRolesModel{}, err
	}

	return data, nil
}

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

func CreateOneService(body HasRolesRequest) lib.ResponseData {
	// validate user input
	if errValidate := inits.MyValidate(body); errValidate != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusBadRequest, Message: errValidate.Error()})
	}

	user, errUsr := users.GetUserOne(body.UserId, "")

	if errUsr != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errUsr.Error()})
	}

	role, errRole := roles.GetRoleOne(body.RoleId, "")

	if errRole != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errRole.Error()})
	}

	_, errCheck := findByUserAndRole(user.ID, role.ID)

	if errCheck == nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: "User already has this role"})
	}

	dataCreated, err := inits.Prisma.HasRoles.CreateOne(
		db.HasRoles.User.Link(
			db.User.ID.Equals(body.UserId),
		),
		db.HasRoles.Role.Link(
			db.Roles.ID.Equals(body.RoleId),
		),
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
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	response := HasRolesWithRelations{
		ID:        dataCreated.ID,
		RoleId:    dataCreated.RoleID,
		UserId:    dataCreated.UserID,
		CreatedAt: dataCreated.CreatedAt,
		UpdatedAt: dataCreated.UpdatedAt,
		Role: RolesInResponse{
			ID:   dataCreated.Role().ID,
			Name: dataCreated.Role().Name,
		},
		User: UserInResponse{
			ID:    dataCreated.User().ID,
			Name:  dataCreated.User().Name,
			Email: dataCreated.User().Email,
		},
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}
