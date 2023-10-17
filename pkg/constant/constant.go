package constant

// user type category
const (
	UserTypeRegular = "regular_user"
	UserTypePremium = "premium_user"
)

const (
	UserInactive = iota
	UserActive
)

var MapStatus = map[int]bool{
	UserInactive: false,
	UserActive:   true,
}
