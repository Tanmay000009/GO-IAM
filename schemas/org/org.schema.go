package orgSchema

import (
	"balkantask/model"
	constants "balkantask/utils"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type OrgResponse struct {
	ID            uuid.UUID               `json:"id,omitempty"`
	Username      string                  `json:"username,omitempty"`
	Email         string                  `json:"email,omitempty"`
	AccountStatus constants.AccountStatus `json:"account_status,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}
type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type SignupInput struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

func MapOrgRecord(user *model.Org) OrgResponse {
	return OrgResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		CreatedAt:     *user.CreatedAt,
		UpdatedAt:     *user.UpdatedAt,
		AccountStatus: user.AccountStatus,
	}
}

// ValidatePassword checks if the password meets complexity requirements.
func (s *SignupInput) ValidatePassword() bool {
	if len(s.Password) < 8 {
		return false
	}

	// Check if the password contains at least one uppercase letter
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(s.Password)

	// Check if the password contains at least one lowercase letter
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(s.Password)

	// Check if the password contains at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(s.Password)

	// Check if the password contains at least one special character
	hasSpecialChar := regexp.MustCompile(`[@$!%*#?&]`).MatchString(s.Password)

	return hasUppercase && hasLowercase && hasDigit && hasSpecialChar
}
