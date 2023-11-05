package constant

// user type category
const (
	UserTypeRegular = "regular_user"
	UserTypePremium = "premium_user"
)

const (
	UserInactive = "0"
	UserActive   = "1"
)

var MapStatus = map[string]bool{
	UserInactive: false,
	UserActive:   true,
}
