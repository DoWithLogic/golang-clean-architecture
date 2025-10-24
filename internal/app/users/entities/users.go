package entities

import (
	"time"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/security"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	ID           int64              `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string             `gorm:"column:name"`
	ContactType  types.CONTACT_TYPE `gorm:"column:contact_type"`
	ContactValue string             `gorm:"column:contact_value"`
	BirthDate    *string            `gorm:"column:birth_date"`
	Language     *types.LANGUAGE    `gorm:"column:language"`
	Password     string             `gorm:"column:password"`
	Status       types.USER_STATUS  `gorm:"column:status"`
	CreatedAt    time.Time          `gorm:"column:created_at"`
	UpdatedAt    *time.Time         `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt     `gorm:"column:deleted_at;index"`
}

func (u User) IsPasswordValid(password string) bool { return u.Password == password }
func (User) TableName() string                      { return types.TABLE_NAME_USERS.String() }

func (u User) ToJWTData(expiresAt int64) security.PayloadToken {
	return security.PayloadToken{
		Data: &security.Data{
			ID:           u.ID,
			ContactType:  u.ContactType,
			ContactValue: u.ContactValue,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
}
