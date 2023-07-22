package roleSchema

import "github.com/google/uuid"

type AddOrDeleteRole struct {
	RoleName string `json:"roleName" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

type TestRole struct {
	RoleName string    `json:"roleName"`
	RoleId   uuid.UUID `json:"roleId"`
}
