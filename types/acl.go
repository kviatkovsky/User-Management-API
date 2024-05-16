package types

import "gorm.io/gorm"

type ACL struct {
	gorm.Model
	UserID uint `gorm:"text;not null" json:"user_id"`
	RoleID uint `gorm:"text;not null" json:"role_id"`
}
