package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;"`
	Username  string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	Role      *string    `gorm:"type:varchar(50);default:'user';not null"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.New()
	return
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
