package users

import (
	"time"
)

type UserRequest struct {
	Name     string `json:"name" validate:"min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=8,max=20"`
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
