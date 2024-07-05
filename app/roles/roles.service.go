package roles

import (
	"context"
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/inits"
	"golang-fiber-prisma/lib"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getOne(id int, name string) (*db.RolesModel, error) {
	var whereUnique db.RolesEqualsUniqueWhereParam

	if name == "" {
		whereUnique = db.Roles.ID.Equals(id)
	} else {
		whereUnique = db.Roles.Name.Equals(name)
	}

	role, err := inits.Prisma.Roles.FindUnique(whereUnique).Exec(context.Background())

	if err != nil {
		return &db.RolesModel{}, err
	}

	if err := lib.CheckDeletedRecord(role.DeletedAt()); err != nil {
		return &db.RolesModel{}, err
	}

	return role, nil
}

func GetAllRolesService(query RoleQueryRequest) lib.ResponseData {
	offset := (query.Page - 1) * query.PerPage
	roles, err := inits.Prisma.Roles.FindMany(
		db.Roles.Or(
			db.Roles.DeletedAt.IsNull(),
			db.Roles.DeletedAt.Equals(time.Time{}),
		),
		db.Roles.Name.Contains(query.Name),
	).OrderBy(
		db.Roles.Name.Order(db.ASC),
	).Skip(offset).Take(query.PerPage).Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: roles})
}

func GetRoleByIdService(id int) lib.ResponseData {
	role, err := getOne(id, "")

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: role})
}
