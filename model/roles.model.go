package model

type Role struct {
	BaseModel
	Name string `gorm:"type:varchar(100);not null; uniqueIndex"`
	Type string `gorm:"type:varchar(100);not null"`
}
