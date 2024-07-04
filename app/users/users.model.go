package users

import (
	"context"
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/lib"
	"time"
)

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserQueryRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type UserResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func getOne(id int, email string, prisma *db.PrismaClient) (*db.UserModel, error) {
	var whereUnique db.UserEqualsUniqueWhereParam

	if email == "" {
		whereUnique = db.User.ID.Equals(id)
	} else {
		whereUnique = db.User.Email.Equals(email)
	}

	user, err := prisma.User.FindUnique(whereUnique).Exec(context.Background())

	if err != nil {
		return &db.UserModel{}, err
	}

	if err := lib.CheckDeletedRecord(user.DeletedAt()); err != nil {
		return &db.UserModel{}, err
	}

	return user, nil
}

func GetAllUsers(query UserQueryRequest, prisma *db.PrismaClient) ([]UserResponse, error) {
	offset := (query.Page - 1) * query.PerPage
	users, err := prisma.User.FindMany(
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
		return []UserResponse{}, err
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

	return response, nil
}

func GetUserById(id int, prisma *db.PrismaClient) (UserResponse, error) {
	user, err := getOne(id, "", prisma)

	if err != nil {
		return UserResponse{}, err
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func GetUserByEmail(email string, prisma *db.PrismaClient) (UserResponse, error) {
	user, err := getOne(0, email, prisma)

	if err != nil {
		return UserResponse{}, err
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func CreateOne(user UserRequest, prisma *db.PrismaClient) (UserResponse, error) {
	newUser, err := prisma.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(user.Password),
	).Exec(context.Background())

	if err != nil {
		return UserResponse{}, err
	}

	response := UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return response, nil
}

func UpdateOne(id int, data UserRequest, prisma *db.PrismaClient) (UserResponse, error) {
	_, errCheck := getOne(id, "", prisma)

	if errCheck != nil {
		return UserResponse{}, errCheck
	}

	user, err := prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.Name.Set(data.Name),
		db.User.Email.Set(data.Email),
		db.User.Password.Set(data.Password),
	).Exec(context.Background())

	if err != nil {
		return UserResponse{}, err
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func DeleteOne(id int, prisma *db.PrismaClient) (UserResponse, error) {
	_, errCheck := getOne(id, "", prisma)

	if errCheck != nil {
		return UserResponse{}, errCheck
	}

	user, err := prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		db.User.DeletedAt.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		return UserResponse{}, err
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}
