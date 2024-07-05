package roles

import "time"

type RoleQueryRequest struct {
	Name    string `json:"name"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type RoleResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleRequest struct {
	Name string `json:"name" validate:"required"`
}
