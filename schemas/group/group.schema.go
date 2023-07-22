package groupSchema

import "github.com/google/uuid"

type AddOrDeleteRole struct {
	RoleId    uuid.UUID `json:"roleId"`
	RoleName  string    `json:"roleName"`
	GroupId   uuid.UUID `json:"groupId"`
	GroupName string    `json:"groupName"`
}

type CreateGroup struct {
	Name      string      `json:"name" validate:"required"`
	RoleIds   []uuid.UUID `json:"roleIds"`
	RoleNames []string    `json:"roleNames"`
}
