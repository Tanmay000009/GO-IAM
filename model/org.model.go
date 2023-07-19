package model

import (
	"time"

	"github.com/google/uuid"
)

type Org struct {
	BaseModel
	ID        uuid.UUID  `gorm:"type:uuid; not null;"`
	Username  string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}
