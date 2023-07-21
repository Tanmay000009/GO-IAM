package orgrepository

import (
	"balkantask/database"
	"balkantask/model"
	constants "balkantask/utils"
)

func FindOrgs() ([]model.Org, error) {
	var orgs []model.Org
	db := database.DB
	err := db.Find(&orgs).Where("account_status != ?", constants.DELETED).Error
	return orgs, err
}

func FindOrgById(id string) (model.Org, error) {
	var org model.Org
	db := database.DB
	err := db.First(&org, "id = ?", id).Error
	return org, err
}

func FindOrgByEmail(email string) (model.Org, error) {
	var org model.Org
	db := database.DB
	err := db.First(&org, "email = ?", email).Error
	return org, err
}

func CreateOrg(org model.Org) (model.Org, error) {
	db := database.DB
	err := db.Create(&org).Error
	return org, err
}

func UpdateOrg(org model.Org) (model.Org, error) {
	db := database.DB
	err := db.Save(&org).Error
	return org, err
}

func DeleteOrg(org model.Org) (model.Org, error) {
	db := database.DB
	err := db.Delete(&org).Error
	return org, err
}
