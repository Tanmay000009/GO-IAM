package model

import "github.com/google/uuid"

type Role struct {
	ID    uuid.UUID `gorm:"type:uuid; not null;"`
	Name  string    `gorm:"type:varchar(100);not null"`
	Users []User    `gorm:"many2many:user_roles;"`
}
