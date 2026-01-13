package user

import (
	"database/sql/driver"
	"fmt"
	"rest-fiber/internal/enums"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role enums.EUserRoleType

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name            string         `gorm:"type:varchar(255);not null;column:name"`
	Email           string         `gorm:"type:varchar(255);unique;not null;column:email"`
	AvatarURL       string         `gorm:"type:text;null;default:null;column:avatar_url"`
	Password        string         `gorm:"type:varchar(255);not null;column:password"`
	IsEmailVerified bool           `gorm:"type:boolean;not null;column:is_email_verified;default:false"`
	Role            Role           `gorm:"type:role_type;not null;default:'MEMBER'"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}

func (r *Role) Scan(value any) error {
	if value == nil {
		*r = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*r = Role(string(v))
	case string:
		*r = Role(v)
	default:
		return fmt.Errorf("cannot scan %T into Role", value)
	}
	return nil
}
func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}
