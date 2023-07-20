package model

type Role struct {
	ID    string `gorm:"type:uuid; not null;"`
	Name  string `gorm:"type:varchar(100);not null"`
	Users []User `gorm:"many2many:user_roles;"`
}
