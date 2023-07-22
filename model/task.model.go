package model

type Task struct {
	BaseModel
	Name  string `gorm:"type:varchar(100);not null; uniqueIndex"`
	Roles []Role `gorm:"many2many:task_roles;constraint:OnDelete:CASCADE;"`
}

func (Task) PrimaryKey() string {
	return "Id"
}
