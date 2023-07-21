package groupSchema

import "github.com/google/uuid"

type AddOrDeleteRole struct {
	RoleId  uuid.UUID `json:"roleId" validate:"required"`
	GroupId uuid.UUID `json:"groupId" validate:"required"`
}

type CreateGroup struct {
	Name    string      `json:"name" validate:"required"`
	RoleIds []uuid.UUID `json:"roleIds" validate:"required"`
}
