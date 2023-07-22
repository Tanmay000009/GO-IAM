package roles

import (
	"balkantask/model"

	"github.com/google/uuid"
)

type Role string

const (
	OrgFullAccess    Role = "ORG_FULL_ACCESS"
	UserFullAccess   Role = "USER_FULL_ACCESS"
	UserWriteAccess  Role = "USER_WRITE_ACCESS"
	UserReadAccess   Role = "USER_READ_ACCESS"
	RoleFullAccess   Role = "ROLE_FULL_ACCESS"
	RoleWriteAccess  Role = "ROLE_WRITE_ACCESS"
	RoleReadAccess   Role = "ROLE_READ_ACCESS"
	OrgWriteAccess   Role = "ORG_WRITE_ACCESS"
	OrgReadAccess    Role = "ORG_READ_ACCESS"
	GroupReadAccess  Role = "GROUP_READ_ACCESS"
	GroupWriteAccess Role = "GROUP_WRITE_ACCESS"
	GroupFullAccess  Role = "GROUP_FULL_ACCESS"
	TasksReadAccess  Role = "TASKS_READ_ACCESS"
	TasksWriteAccess Role = "TASKS_WRITE_ACCESS"
	TasksFullAccess  Role = "TASKS_FULL_ACCESS"
)

func UserIsAuthorized(roles []model.Role, group []model.Group, targetRoles []Role) bool {

	userRoles := []model.Role{}
	for _, role := range roles {
		userRoles = append(userRoles, role)
	}
	for _, group := range group {
		for _, role := range group.Roles {
			userRoles = append(userRoles, role)
		}
	}

	uniqueRoles := RemoveDuplicates(userRoles)

	for _, targetRole := range targetRoles {
		for _, role := range uniqueRoles {
			if Role(role.Name) == targetRole {
				return true
			}
		}
	}
	return false
}

func UserHasRole(roles []model.Role, targetRole model.Role) bool {
	for _, role := range roles {
		if role.ID == targetRole.ID {
			return true
		}
	}
	return false
}

func GroupHasRole(roles []model.Role, targetRoles []model.Role) bool {
	for _, targetRole := range targetRoles {
		for _, role := range roles {
			if role.ID == targetRole.ID {
				return true
			}
		}
	}
	return false
}

func TaskHasRole(roles []model.Role, targetRoles []model.Role) bool {
	for _, targetRole := range targetRoles {
		for _, role := range roles {
			if role.ID == targetRole.ID {
				return true
			}
		}
	}
	return false
}

func UserHasGroup(groups []model.Group, targetGroups []model.Group) bool {
	for _, targetGroup := range targetGroups {
		for _, group := range groups {
			if group.ID == targetGroup.ID {
				return true
			}
		}
	}
	return false
}

func UserHasTaskAuthorization(roles []model.Role, group []model.Group, targetRoles []model.Role) bool {
	userRoles := []model.Role{}
	for _, role := range roles {
		userRoles = append(userRoles, role)
	}
	for _, group := range group {
		for _, role := range group.Roles {
			userRoles = append(userRoles, role)
		}
	}

	uniqueRoles := RemoveDuplicates(userRoles)

	for _, targetRole := range targetRoles {
		for _, task := range uniqueRoles {
			if task.ID == targetRole.ID {
				return true
			}
		}
	}
	return false
}

func RemoveDuplicates(rolesExist []model.Role) []model.Role {

	uniqueRolesMap := make(map[uuid.UUID]struct{})
	var uniqueRoles []model.Role

	for _, role := range rolesExist {
		if _, found := uniqueRolesMap[role.ID]; !found {
			uniqueRolesMap[role.ID] = struct{}{}
			uniqueRoles = append(uniqueRoles, role)
		}
	}

	return uniqueRoles
}
