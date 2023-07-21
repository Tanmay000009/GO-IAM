package model

import (
	constants "balkantask/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Username      string                  `gorm:"type:varchar(100);not null; uniqueIndex"`
	Password      string                  `gorm:"type:varchar(100);not null"`
	OrgID         uuid.UUID               `gorm:"type:uuid;"`
	Roles         []Role                  `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Org           *Org                    `gorm:"foreignKey:OrgID;constraint:OnDelete:CASCADE;"`
	AccountStatus constants.AccountStatus `gorm:"type:varchar(100);not null;default:'active'"`
	CreatedAt     *time.Time              `gorm:"not null;default:now()"`
	UpdatedAt     *time.Time              `gorm:"not null;default:now()"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
