package rolesRepo

import (
	"balkantask/database"
	"balkantask/model"
	roles "balkantask/utils"

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

func GetRoleByName(name string) (model.Role, error) {
	db := database.DB
	var role model.Role
	err := db.First(&role, "name = ?", name).Error

	return role, err
}

func CreateRole(role model.Role) (model.Role, error) {
	db := database.DB
	err := db.Create(&role).Error

	return role, err
}

func DeleteRoleById(id uuid.UUID) error {
	db := database.DB

	var role roles.Role
	err := db.Delete(&role, "id = ?", id).Error

	return err
}
