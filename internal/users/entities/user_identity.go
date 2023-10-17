package entities

type Identity struct {
	Email  string `json:"email,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
}
