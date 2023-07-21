package groupRepo

import (
	"balkantask/database"
	"balkantask/model"

	"github.com/google/uuid"
)

func GetAllGroups() ([]model.Group, error) {
	var groups []model.Group
	db := database.DB
	err := db.Preload("Roles").Find(&groups).Error
	return groups, err
}

func GetGroupById(id uuid.UUID) (model.Group, error) {
	var group model.Group
	db := database.DB
	err := db.Preload("Roles").Where("id = ?", id).First(&group).Error
	return group, err
}

func CreateGroup(group *model.Group) (*model.Group, error) {
	db := database.DB
	err := db.Create(&group).Error
	return group, err
}

func DeleteGroup(group *model.Group) error {
	db := database.DB
	err := db.Model(&group).Association("Roles").Clear()
	err = db.Model(&group).Association("Users").Clear()
	err = db.Delete(&group).Error
	return err
}

func AddRoleToGroup(group model.Group, role model.Role) (model.Group, error) {
	db := database.DB
	err := db.Model(&group).Association("Roles").Append(&role)
	return group, err

}

func RemoveRoleFromGroup(group model.Group, role model.Role) (model.Group, error) {
	db := database.DB
	err := db.Model(&group).Association("Roles").Delete(&role)
	return group, err
}
