package enums

type ETokenType string

const (
	TokenAccess  ETokenType = "access"
	TokenRefresh ETokenType = "refresh"
	TokenVerify  ETokenType = "verify"
)
