package userSchema

import (
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	constants "balkantask/utils"
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password,omitempty" validate:"omitempty,min=8"`
	PasswordConfirm string `json:"passwordConfirm,omitempty" validate:"omitempty,min=8"`
}

type UserResponse struct {
	ID            uuid.UUID               `json:"id,omitempty"`
	Username      string                  `json:"username,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Roles         []model.Role            `json:"roles"`
	OrgId         uuid.UUID               `json:"org_id,omitempty"`
	AccountStatus constants.AccountStatus `json:"account_status,omitempty"`
}

type UserResponseWithOrg struct {
	ID            uuid.UUID               `json:"id,omitempty"`
	Username      string                  `json:"username,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Roles         []model.Role            `json:"roles"`
	OrgId         uuid.UUID               `json:"org_id,omitempty"`
	Org           orgSchema.OrgResponse   `json:"org,omitempty"`
	AccountStatus constants.AccountStatus `json:"account_status,omitempty"`
}

type CreateUserResponse struct {
	ID            uuid.UUID               `json:"id,omitempty"`
	Username      string                  `json:"username,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Roles         []model.Role            `json:"roles"`
	OrgId         uuid.UUID               `json:"org_id,omitempty"`
	AccountStatus constants.AccountStatus `json:"account_status,omitempty"`
	Passcode      string                  `json:"passcode,omitempty"`
}

type AddOrDeleteRole struct {
	RoleId uuid.UUID `json:"role_id" validate:"required"`
	UserId uuid.UUID `json:"user_id" validate:"required"`
}

func MapUserRecord(user *model.User) UserResponse {

	if user == nil || user.ID == uuid.Nil {

		return UserResponse{
			ID: uuid.Nil,
		}
	}

	return UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		CreatedAt:     *user.CreatedAt,
		UpdatedAt:     *user.UpdatedAt,
		Roles:         user.Roles,
		OrgId:         user.OrgID,
		AccountStatus: user.AccountStatus,
	}
}

func MapUserRecordWithOrg(user *model.User) UserResponseWithOrg {
	return UserResponseWithOrg{
		ID:            user.ID,
		Username:      user.Username,
		CreatedAt:     *user.CreatedAt,
		UpdatedAt:     *user.UpdatedAt,
		Roles:         user.Roles,
		OrgId:         user.OrgID,
		Org:           orgSchema.MapOrgRecord(user.Org),
		AccountStatus: user.AccountStatus,
	}
}

type SignInInput struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
