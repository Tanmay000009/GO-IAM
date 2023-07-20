package model

type Role struct {
	BaseModel
	Name  string `gorm:"type:varchar(100);not null"`
	Users []User `gorm:"many2many:user_roles;"`
	Type  string `gorm:"type:varchar(100);not null"`
}
