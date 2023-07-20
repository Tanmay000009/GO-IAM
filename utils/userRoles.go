package roles

import "balkantask/model"

type Role string

const (
	OrgFullAccess   Role = "ORG_FULL_ACCESS"
	UserFullAccess  Role = "USER_FULL_ACCESS"
	UserWriteAccess Role = "USER_WRITE_ACCESS"
	UserReadAccess  Role = "USER_READ_ACCESS"
	RoleFullAccess  Role = "ROLE_FULL_ACCESS"
	RoleWriteAccess Role = "ROLE_WRITE_ACCESS"
	RoleReadAccess  Role = "ROLE_READ_ACCESS"
	OrgWriteAccess  Role = "ORG_WRITE_ACCESS"
	OrgReadAccess   Role = "ORG_READ_ACCESS"
)

func HasAnyRole(roles []model.Role, targetRoles ...Role) bool {
	for _, targetRole := range targetRoles {
		for _, role := range roles {
			if Role(role.Name) == targetRole {
				return true
			}
		}
	}
	return false
}
