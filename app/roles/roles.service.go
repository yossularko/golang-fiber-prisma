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
		db.Roles.DeletedAt.IsNull(),
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

func CreateOneService(body RoleRequest) lib.ResponseData {
	// validate user input
	if err := inits.MyValidate(body); err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusBadRequest, Message: err.Error()})
	}

	// check if role exist
	_, errRoleCheck := getOne(0, body.Name)

	if errRoleCheck == nil {
		message := "Role already exist"
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: &message})
	}

	newRole, err := inits.Prisma.Roles.CreateOne(
		db.Roles.Name.Set(body.Name),
	).Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: newRole})
}

func UpdateOneService(id int, body RoleRequest) lib.ResponseData {
	// validate user input
	if err := inits.MyValidate(body); err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusBadRequest, Message: err.Error()})
	}

	// check role
	_, errCheck := getOne(id, "")

	if errCheck != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errCheck.Error()})
	}

	newRole, err := inits.Prisma.Roles.
		FindUnique(db.Roles.ID.Equals(id)).
		Update(db.Roles.Name.Set(body.Name)).
		Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: newRole})
}

func DeleteOneService(id int) lib.ResponseData {
	// check role
	_, errCheck := getOne(id, "")

	if errCheck != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errCheck.Error()})
	}

	newRole, err := inits.Prisma.Roles.
		FindUnique(db.Roles.ID.Equals(id)).
		Update(db.Roles.DeletedAt.Set(time.Now())).
		Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: newRole})
}
