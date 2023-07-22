package rolesRepo

import (
	"balkantask/database"
	"balkantask/model"

	"github.com/google/uuid"
)

func GetAllRoles() ([]model.Role, error) {
	db := database.DB
	var roles []model.Role
	err := db.Find(&roles).Error

	return roles, err
}

func GetRoleById(id uuid.UUID) (model.Role, error) {
	db := database.DB
	var role model.Role
	err := db.First(&role, "id = ?", id).Error

	return role, err
}

func GetRolesByIds(ids []uuid.UUID) ([]model.Role, error) {
	db := database.DB
	var roles []model.Role
	err := db.Find(&roles, "id IN ?", ids).Error
	return roles, err
}

func GetRoleByName(name string) (model.Role, error) {
	db := database.DB
	var role model.Role
	err := db.First(&role, "name = ?", name).Error

	return role, err
}

func GetRolesByNames(name []string) ([]model.Role, error) {
	db := database.DB
	var role []model.Role
	err := db.Find(&role, "name IN ?", name).Error

	return role, err
}

func CreateRole(role model.Role) (model.Role, error) {
	db := database.DB
	err := db.Create(&role).Error

	return role, err
}

func DeleteRole(role *model.Role) error {
	db := database.DB

	err := db.Model(&role).Association("Users").Clear()
	err = db.Model(&role).Association("Groups").Clear()
	err = db.Delete(&role).Error
	return err
}
