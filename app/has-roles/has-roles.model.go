package hasroles

import "time"

type HasRolesRequest struct {
	UserId int `json:"user_id" validate:"required"`
	RoleId int `json:"role_id" validate:"required"`
}

type RolesInResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type UserInResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type HasRolesWithRelations struct {
	ID        int             `json:"id"`
	RoleId    int             `json:"role_id"`
	UserId    int             `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Role      RolesInResponse `json:"role"`
	User      UserInResponse  `json:"user"`
}

type HasRolesWithRole struct {
	ID        int             `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Roles     RolesInResponse `json:"role"`
}
