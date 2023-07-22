package roleSchema

import "github.com/google/uuid"

type AddOrDeleteRole struct {
	RoleName string `json:"roleName" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

type CreateGroup struct {
	Name      string      `json:"name" validate:"required"`
	RoleIds   []uuid.UUID `json:"roleIds"`
	RoleNames []string    `json:"roleNames"`
}
