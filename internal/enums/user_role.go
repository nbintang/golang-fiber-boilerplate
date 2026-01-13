package enums

type EUserRoleType string

const (
	Admin  EUserRoleType = "ADMIN"
	Member EUserRoleType = "MEMBER"
)

func (r EUserRoleType) IsValid() bool {
	return r == Admin || r == Member
}

func (r EUserRoleType) IsAdmin() bool {
	return r == Admin
}

func (r EUserRoleType) IsMember() bool {
	return r == Member
}
