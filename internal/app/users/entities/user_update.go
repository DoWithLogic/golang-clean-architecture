package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
)

type UpdateUser struct {
	ID           int64               `gorm:"column:id;primaryKey;autoIncrement"`
	Name         *string             `gorm:"column:name"`
	ContactType  *types.CONTACT_TYPE `gorm:"column:contact_type"`
	ContactValue *string             `gorm:"column:contact_value"`
	BirthDate    *string             `gorm:"column:birth_date"`
	Language     *types.LANGUAGE     `gorm:"column:language"`
	Password     *string             `gorm:"column:password"`
	Status       *types.USER_STATUS  `gorm:"column:status"`
	UpdatedAt    time.Time           `gorm:"column:updated_at"`
}
