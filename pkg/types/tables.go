package types

type TABLE_NAME string

const (
	TABLE_NAME_USERS TABLE_NAME = "users"
)

func (t TABLE_NAME) String() string { return string(t) }
