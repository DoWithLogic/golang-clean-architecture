package entities

import (
	"time"

	jwtPkg "github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/golang-jwt/jwt/v5"
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
func (User) TableName() string                      { return "users" }

func (u User) ToJWTData(expiresAt int64) *jwtPkg.JWTClaims {
	expiredTime := time.Now().Add(time.Minute * time.Duration(expiresAt))

	return &jwtPkg.JWTClaims{
		Data: &jwtPkg.Data{
			ID:           u.ID,
			ContactType:  u.ContactType,
			ContactValue: u.ContactValue,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
}
