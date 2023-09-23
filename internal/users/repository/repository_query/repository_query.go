package repository_query

import _ "embed"

var (
	//go:embed users/insert.sql
	InsertUsers string
	//go:embed users/update.sql
	UpdateUsers string
	//go:embed users/select.sql
	GetUserByID string
)
