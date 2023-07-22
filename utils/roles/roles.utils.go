package roles

import "balkantask/model"

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
)

func HasAnyRole(roles []model.Role, group []model.Group, targetRoles []Role) bool {
	userRoles := append(roles, group[0].Roles...)
	for _, targetRole := range targetRoles {
		for _, role := range userRoles {
			if Role(role.Name) == targetRole {
				return true
			}
		}
	}
	return false
}

func UserHasRole(roles []model.Role, targetRoles []Role) bool {
	for _, targetRole := range targetRoles {
		for _, role := range roles {
			if Role(role.Name) == targetRole {
				return true
			}
		}
	}
	return false
}

func GroupHasRole(roles []model.Role, targetRoles []Role) bool {
	for _, targetRole := range targetRoles {
		for _, role := range roles {
			if Role(role.Name) == targetRole {
				return true
			}
		}
	}
	return false
}

func HasAnyGroup(groups []model.Group, targetGroups []model.Group) bool {
	for _, targetGroup := range targetGroups {
		for _, group := range groups {
			if group.ID == targetGroup.ID {
				return true
			}
		}
	}
	return false
}
