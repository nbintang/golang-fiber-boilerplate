package token

import "rest-fiber/internal/enums"

type GenerateTokenParams struct {
	ID    string
	Email string
	Role  enums.EUserRoleType
	JTI   string
	Type  enums.ETokenType
}