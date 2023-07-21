package model

import (
	constants "balkantask/utils"
	"time"
)

type Org struct {
	BaseModel
	Username      string                  `gorm:"type:varchar(100);not null"`
	Email         string                  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password      string                  `gorm:"type:varchar(100);not null"`
	Users         []User                  `gorm:"foreignKey:OrgID;constraint:OnDelete:CASCADE;"`
	AccountStatus constants.AccountStatus `gorm:"type:varchar(100);not null;default:'active'"`
	CreatedAt     *time.Time              `gorm:"not null;default:now()"`
	UpdatedAt     *time.Time              `gorm:"not null;default:now()"`
}
