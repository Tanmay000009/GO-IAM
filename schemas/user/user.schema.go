package userSchema

import (
	"balkantask/model"
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []string  `json:"roles"`
}

func MapUserRecord(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
		Roles:     user.Roles,
	}
}

type SignInInput struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
