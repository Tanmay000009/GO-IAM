package model

type Role struct {
	BaseModel
	Name   string  `gorm:"type:varchar(100);not null; uniqueIndex"`
	Type   string  `gorm:"type:varchar(100);not null"`
	Users  []User  `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
	Groups []Group `gorm:"many2many:group_roles;constraint:OnDelete:CASCADE;"`
}
