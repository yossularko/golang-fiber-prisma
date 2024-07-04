package users

import (
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/lib"

	"github.com/gofiber/fiber/v2"
)

func validateStoreRequest(user UserRequest) []lib.ValidationResponse {
	// Define validation rules
	rules := lib.ValidationRules{
		"Name": func(value interface{}) bool {
			// Name must be a string and not empty
			name, ok := value.(string)
			return ok && name != ""
		},
		"Email": func(value interface{}) bool {
			// Email must be a string and not empty and must be a valid email
			email, ok := value.(string)
			return ok && email != "" && lib.ValidateEmail(email)
		},
		"Password": func(value interface{}) bool {
			// Password must be a string and not empty and must be at least 8 characters long
			password, ok := value.(string)
			return ok && password != "" && len(password) >= 8
		},
	}

	// Convert UserRequest to map
	userMap := map[string]interface{}{
		"Name":     user.Name,
		"Email":    user.Email,
		"Password": user.Password,
	}

	// Validate user input
	errors := lib.ValidateRequest(userMap, rules)

	return errors
}

func GetAllUsersService(query UserQueryRequest, prisma *db.PrismaClient) lib.ResponseData {
	users, err := GetAllUsers(query, prisma)

	if err != nil {
		lib.ResponseError(lib.ResponseProps{Code: fiber.StatusInternalServerError, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusOK, Data: users})
}

func CreateOneService(user UserRequest, prisma *db.PrismaClient) lib.ResponseData {
	// validate user input
	errors := validateStoreRequest(user)
	if len(errors) > 0 {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusBadRequest, Message: &errors})
	}
	// check if email already exist
	_, errUsrCekEmail := GetUserByEmail(user.Email, prisma)

	if errUsrCekEmail == nil {
		message := "Email already exist"
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: &message})
	}

	hashPasswrd, _ := lib.HashPassword(user.Password)
	newUser, err := CreateOne(UserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashPasswrd,
	}, prisma)

	if err != nil {
		return lib.ResponseError(lib.ResponseProps{Code: fiber.StatusConflict, Message: err.Error()})
	}

	return lib.ResponseSuccess(lib.ResponseProps{Code: fiber.StatusCreated, Data: newUser})
}
