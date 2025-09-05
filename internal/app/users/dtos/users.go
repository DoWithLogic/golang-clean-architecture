package dtos

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
)

type User struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	ContactType  types.CONTACT_TYPE `json:"contact_type"`
	ContactValue string             `json:"contact_value"`
	BirthDate    *string            `json:"birth_date"`
	Language     *types.LANGUAGE    `json:"language"`
	Password     string             `json:"password"`
	Status       types.USER_STATUS  `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at"`
}
