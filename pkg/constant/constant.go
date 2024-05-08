package constant

// user type category
const (
	UserTypeRegular = "regular_user"
	UserTypePremium = "premium_user"
)

const (
	UserInactive = "inactive"
	UserActive   = "active"
)

var MapStatus = map[string]bool{
	UserInactive: false,
	UserActive:   true,
}

const (
	AuthorizationHeaderKey = "Authorization"
	AuthCredentialKey      = "authCredential"
)
