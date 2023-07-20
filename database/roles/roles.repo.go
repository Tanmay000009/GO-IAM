package rolesRepo

import (
	"balkantask/database"
	"balkantask/model"
	roles "balkantask/utils"

	"github.com/google/uuid"
)

func GetAllRoles() ([]roles.Role, error) {
	db := database.DB
	var roles []roles.Role
	err := db.Find(&roles).Error

	return roles, err
}

func GetRoleById(id uuid.UUID) (roles.Role, error) {
	db := database.DB
	var role roles.Role
	err := db.First(&role, "id = ?", id).Error

	return role, err
}

func GetRoleByName(name string) (roles.Role, error) {
	db := database.DB
	var role roles.Role
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
