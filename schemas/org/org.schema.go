package orgSchema

import (
	"balkantask/model"
	"time"

	"github.com/google/uuid"
)

type OrgReponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []string  `json:"roles,omitempty"`
}

func MapOrgRecord(user *model.Org) OrgReponse {
	return OrgReponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
		Roles:     user.Roles,
	}
}
