package userrepository

import (
	"balkantask/database"
	"balkantask/model"
)

func FindUsers() ([]model.User, error) {
	var users []model.User
	db := database.DB
	err := db.Find(&users).Error
	return users, err
}

func FindUserById(id string) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.First(&user, "id = ?", id).Error
	return user, err
}

func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	db := database.DB
	err := db.First(&user, "username = ?", username).Error
	return user, err
}

func CreateUser(user model.User) (model.User, error) {
	db := database.DB
	err := db.Create(&user).Error
	return user, err
}

func UpdateUser(user model.User) (model.User, error) {
	db := database.DB
	err := db.Save(&user).Error
	return user, err
}

func DeleteUser(user model.User) (model.User, error) {
	db := database.DB
	err := db.Delete(&user).Error
	return user, err
}
