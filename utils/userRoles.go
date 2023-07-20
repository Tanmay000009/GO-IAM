package roles

type Role string

const (
	OrgFullAccess   Role = "ORG_FULL_ACCESS"
	UserFullAccess  Role = "USER_FULL_ACCESS"
	UserWriteAccess Role = "USER_WRITE_ACCESS"
	UserReadAccess  Role = "USER_READ_ACCESS"
)
