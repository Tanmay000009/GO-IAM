package orgrepository

import (
	"balkantask/database"
	"balkantask/model"
	orgSchema "balkantask/schemas/org"
	constants "balkantask/utils"

	"github.com/google/uuid"
)

func FindOrgs() ([]model.Org, error) {
	var orgs []model.Org
	db := database.DB
	err := db.Find(&orgs).Where("account_status != ?", constants.DELETED).Error
	return orgs, err
}

func FindOrgById(id uuid.UUID) (model.Org, error) {
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

func UpdateOrg(org model.Org) (orgSchema.OrgResponse, error) {
	db := database.DB
	err := db.Save(&org).Error

	org_ := orgSchema.MapOrgRecord(&org)

	return org_, err
}

func DeleteOrg(org model.Org) (model.Org, error) {
	db := database.DB
	err := db.Delete(&org).Error
	return org, err
}
