package userrepository

import (
	"balkantask/database"
	"balkantask/model"
	userSchema "balkantask/schemas/user"
	constants "balkantask/utils"
	"time"

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
	err := db.Preload("Roles").Preload("Groups").Preload("Org").Where("id = ? AND account_status != ?", id, constants.DELETED).First(&user).Error
	return user, err
}

func FindUserById(id uuid.UUID) (userSchema.UserResponse, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").Preload("Groups").Where("id = ? AND account_status != ?", id, constants.DELETED).First(&user).Error
	user_ := userSchema.MapUserRecord(&user)

	return user_, err
}

func FindUserWithOrgById(id string) (userSchema.UserResponseWithOrg, error) {
	var user model.User

	db := database.DB
	err := db.Preload("Roles").Preload("Groups").Preload("Orgs").Where("id = ? AND account_status != ?", id, constants.DELETED).First(&user).Error
	if err != nil {
		return userSchema.UserResponseWithOrg{}, err
	}
	if user.ID == uuid.Nil {
		return userSchema.UserResponseWithOrg{}, err
	}
	userWithOrg := userSchema.MapUserRecordWithOrg(&user)

	return userWithOrg, err
}

func FindUserWithOrgByUsername(username string) (userSchema.UserResponseWithOrg, error) {
	var user model.User

	db := database.DB

	// Perform a left join on orgs table using Preload
	err := db.Preload("Roles").Preload("Groups").Preload("Org").Where("username = ? AND account_status != ?", username, constants.DELETED).First(&user).Error

	userWithOrg := userSchema.MapUserRecordWithOrg(&user)

	return userWithOrg, err
}

func FindUserByUsernameWithPassword(username string) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").Preload("Groups").Where("username = ?", username).First(&user).Error
	return user, err
}

func FindUserByOrgAndUsernameWithPassword(username string, orgId string) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.Preload("Roles").Preload("Groups").Where("username = ? AND org_id = ?", username, orgId).First(&user).Error
	return user, err
}

func FindUserByUsername(username string) (*model.User, error) {
	var user *model.User
	db := database.DB
	err := db.Where("username = ? AND account_status != ?", username, constants.DELETED).First(&user).Error

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
	err := db.Where("account_status != ? AND org_id = ? AND roles LIKE ?", constants.DELETED, orgId, "%"+role+"%").Find(&users).Error

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
	err := db.Model(&user).Association("Roles").Clear()
	err = db.Model(&user).Association("Groups").Clear()
	err = db.Delete(&user).Error

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

func AddGroupToUser(group model.Group, user model.User) (model.User, error) {
	db := database.DB
	err := db.Model(&user).Association("Groups").Append(&group)
	return user, err
}

func DeleteGroupFromUser(group model.Group, user model.User) (model.User, error) {
	db := database.DB
	err := db.Model(&user).Association("Groups").Delete(&group)
	return user, err
}

func GetDeactivatedUserForThreshold(threshold time.Time) ([]model.User, error) {
	var users []model.User
	db := database.DB
	err := db.Where("account_status = ? AND updated_at < ?", constants.DEACTIVATED, threshold).Find(&users).Error
	return users, err

}

func DeleteUsers(users []model.User) error {
	db := database.DB
	err := db.Model(&users).Association("Roles").Clear()
	err = db.Model(&users).Association("Groups").Clear()
	err = db.Delete(&users).Error
	return err
}
