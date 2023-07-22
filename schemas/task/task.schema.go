package taskSchema

import "github.com/google/uuid"

type CreateTask struct {
	Name      string      `json:"name" validate:"required"`
	RoleIds   []uuid.UUID `json:"roleIds"`
	RoleNames []string    `json:"roleNames"`
}

type AddOrDeleteRole struct {
	RoleId   uuid.UUID `json:"roleId"`
	RoleName string    `json:"roleName"`
	TaskId   uuid.UUID `json:"taskId"`
	TaskName string    `json:"taskName"`
}

type TestTask struct {
	TaskName string    `json:"taskName"`
	TaskId   uuid.UUID `json:"taskId"`
}
