package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid; not null;"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	base.ID = uuid.New()
	return
}
