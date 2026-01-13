package identity

import "rest-fiber/internal/enums"

type AuthClaims struct {
	ID    string
	Email string
	Role  enums.EUserRoleType
	JTI   string
}