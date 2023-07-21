package userrepository

import (
	"balkantask/database"
	"balkantask/model"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"

	"github.com/google/uuid"
)

func FindUsers() ([]userSchema.UserResponse, error) {
	var users []model.User
	db := database.DB
	err := db.Where("account_status != ?", constants.DELETED).Find(&users).Error

	var users_ []userSchema.UserResponse
	for _, user := range users {
		users_ = append(users_, userSchema.MapUserRecord(&user))
	}

	return users_, err
}

func FindUserByIdWithPassword(id uuid.UUID) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").First(&user).Where("id = ? AND account_status != ?", id, constants.DELETED).Error
	return user, err
}

func FindUserById(id uuid.UUID) (userSchema.UserResponse, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").First(&user).Where("id = ? AND account_status != ?", id, constants.DELETED).Error
	user_ := userSchema.MapUserRecord(&user)

	return user_, err
}

func FindUserWithOrgById(id string) (userSchema.UserResponseWithOrg, error) {
	var user model.User

	db := database.DB
	err := db.Preload("Roles").Preload("Org").First(&user).Where("id = ? AND account_status != ?", id, constants.DELETED).Error

	userWithOrg := userSchema.MapUserRecordWithOrg(&user)

	return userWithOrg, err
}

func FindUserWithOrgByUsername(username string) (userSchema.UserResponseWithOrg, error) {
	var user model.User

	db := database.DB

	// Perform a left join on orgs table using Preload
	err := db.Preload("Roles").Preload("Org").First(&user).Where("username = ? AND account_status != ?", username, constants.DELETED).Error

	userWithOrg := userSchema.MapUserRecordWithOrg(&user)

	return userWithOrg, err
}

func FindUserByUsernameWithPassword(username string) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").First(&user).Where("account_status != ? AND username = ?", constants.DELETED, username).Error
	return user, err
}

func FindUserByUsername(username string) (*model.User, error) {
	var user *model.User
	db := database.DB
	err := db.First(&user).Where("username = ? AND account_status != ?", username, constants.DELETED).Error

	return user, err
}

func FindUsersByOrgId(orgId uuid.UUID) ([]userSchema.UserResponse, error) {
	var users []model.User
	db := database.DB
	err := db.Where("org_id = ? AND account_status != ?", orgId, constants.DELETED).Preload("Roles").Find(&users).Error

	var users_ []userSchema.UserResponse
	for _, user := range users {
		users_ = append(users_, userSchema.MapUserRecord(&user))
	}

	return users_, err
}

func FindUsersByOrgIdAndRole(orgId uuid.UUID, role string) ([]userSchema.UserResponse, error) {
	var users []model.User
	db := database.DB
	err := db.Find(&users).Where("account_status != ? AND org_id = ? AND roles LIKE ?", constants.DELETED, orgId, "%"+role+"%").Error

	var users_ []userSchema.UserResponse
	for _, user := range users {
		users_ = append(users_, userSchema.MapUserRecord(&user))
	}

	return users_, err
}

func CreateUser(user model.User) (userSchema.UserResponse, error) {
	db := database.DB
	err := db.Create(&user).Error
	user_ := userSchema.MapUserRecord(&user)

	return user_, err
}

func UpdateUser(user model.User) (userSchema.UserResponse, error) {
	db := database.DB
	err := db.Save(&user).Error
	user_ := userSchema.MapUserRecord(&user)

	return user_, err
}

func DeleteUser(user model.User) (bool, error) {
	db := database.DB
	err := db.Delete(&user).Error

	return true, err
}

func AddRoleToUser(role model.Role, user model.User) (model.User, error) {
	db := database.DB
	err := db.Model(&user).Association("Roles").Append(&role)
	return user, err
}

func DeleteRoleFromUser(role model.Role, user model.User) (model.User, error) {
	db := database.DB
	err := db.Model(&user).Association("Roles").Delete(&role)
	return user, err
}
