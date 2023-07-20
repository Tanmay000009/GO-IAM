package model

import (
	userroles "balkantask/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID        `gorm:"type:uuid; not null;"`
	CreatedAt *time.Time       `gorm:"not null;default:now()"`
	UpdatedAt *time.Time       `gorm:"not null;default:now()"`
	Roles     []userroles.Role `gorm:"-"` // Use UserRole type and set the database column type to an array of strings
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	base.ID = uuid.New()
	return
}
