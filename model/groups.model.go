package model

type Group struct {
	BaseModel
	Name  string `gorm:"type:varchar(100);not null; uniqueIndex"`
	Roles []Role `gorm:"many2many:group_roles;constraint:OnDelete:CASCADE;"`
	Users []User `gorm:"many2many:user_groups;constraint:OnDelete:CASCADE;"`
}

func (Group) PrimaryKey() string {
	return "Id"
}
