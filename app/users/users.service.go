package users

import (
	"context"
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/inits"
	"golang-fiber-prisma/lib"
	"time"

	"github.com/gofiber/fiber/v2"
)

// func validateStoreRequest(user UserRequest) []lib.ValidationResponse {
// 	// Define validation rules
// 	rules := lib.ValidationRules{
// 		"Name": func(value interface{}) bool {
// 			// Name must be a string and not empty
// 			name, ok := value.(string)
// 			return ok && name != ""
// 		},
// 		"Email": func(value interface{}) bool {
// 			// Email must be a string and not empty and must be a valid email
// 			email, ok := value.(string)
// 			return ok && email != "" && lib.ValidateEmail(email)
// 		},
// 		"Password": func(value interface{}) bool {
// 			// Password must be a string and not empty and must be at least 8 characters long
// 			password, ok := value.(string)
// 			return ok && password != "" && len(password) >= 8
// 		},
// 	}

// 	// Convert UserRequest to map
// 	userMap := map[string]interface{}{
// 		"Name":     user.Name,
// 		"Email":    user.Email,
// 		"Password": user.Password,
// 	}

// 	// Validate user input
// 	errors := lib.ValidateRequest(userMap, rules)

// 	return errors
// }

func getOne(id int, email string) (*db.UserModel, error) {
	var whereUnique db.UserEqualsUniqueWhereParam

	if email == "" {
		whereUnique = db.User.ID.Equals(id)
	} else {
		whereUnique = db.User.Email.Equals(email)
	}

	user, err := inits.Prisma.User.FindUnique(whereUnique).Exec(context.Background())

	if err != nil {
		return &db.UserModel{}, err
	}

	if err := lib.CheckDeletedRecord(user.DeletedAt()); err != nil {
		return &db.UserModel{}, err
	}

	return user, nil
}

func GetAllUsersService(query UserQueryRequest) lib.ResponseData {
	offset := (query.Page - 1) * query.PerPage
	users, err := inits.Prisma.User.FindMany(
		db.User.DeletedAt.IsNull(),
		db.User.Name.Contains(query.Name),
		db.User.Email.Contains(query.Email),
	).OrderBy(
		db.User.CreatedAt.Order(db.DESC),
	).Select(
		db.User.ID.Field(),
		db.User.Name.Field(),
		db.User.Email.Field(),
		db.User.CreatedAt.Field(),
		db.User.UpdatedAt.Field(),
	).Skip(offset).Take(query.PerPage).Exec(context.Background())

	if err != nil {
		lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	response := []UserResponse{}

	for _, user := range users {
		response = append(response, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}

func GetUserByIdService(id int) lib.ResponseData {
	user, err := getOne(id, "")

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: err.Error()})
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}

func CreateOneService(user UserRequest) lib.ResponseData {
	// validate user input
	if err := inits.MyValidate(user); err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusBadRequest, Message: err.Error()})
	}

	// check if email already exist
	_, errUsrCekEmail := getOne(0, user.Email)

	if errUsrCekEmail == nil {
		message := "Email already exist"
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: &message})
	}

	hashPasswrd, _ := lib.HashPassword(user.Password)

	newUser, err := inits.Prisma.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(hashPasswrd),
	).Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: err.Error()})
	}

	response := UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusCreated, Data: response})
}

func UpdateOneService(id int, data UserRequest) lib.ResponseData {
	initUser, errCheck := getOne(id, "")

	if errCheck != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errCheck.Error()})
	}

	if data.Name != "" {
		errs := inits.Validate.Var(data.Name, "min=5,max=30")
		if errs != nil {
			return lib.ResponseError(lib.ResponseProps{
				Code:    fiber.StatusBadRequest,
				Message: inits.ParseVarValidate("name", errs),
			})
		}
		initUser.Name = data.Name
	}

	if data.Email != "" {
		errs := inits.Validate.Var(data.Email, "required,email")
		if errs != nil {
			return lib.ResponseError(lib.ResponseProps{
				Code:    fiber.StatusBadRequest,
				Message: inits.ParseVarValidate("email", errs),
			})
		}
		initUser.Email = data.Email
	}

	if data.Password != "" {
		errs := inits.Validate.Var(data.Password, "min=8,max=20")
		if errs != nil {
			return lib.ResponseError(lib.ResponseProps{
				Code:    fiber.StatusBadRequest,
				Message: inits.ParseVarValidate("password", errs),
			})
		}
		hashPasswrd, _ := lib.HashPassword(data.Password)
		initUser.Password = hashPasswrd
	}

	user, err := inits.Prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.Name.Set(initUser.Name),
		db.User.Email.Set(initUser.Email),
		db.User.Password.Set(initUser.Password),
	).Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}

func DeleteOneService(id int) lib.ResponseData {
	_, errCheck := getOne(id, "")

	if errCheck != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusNotFound, Message: errCheck.Error()})
	}

	user, err := inits.Prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.DeletedAt.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: response})
}
