package types

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	ID          uint   `gorm:"primarykey"`
	RoleLabel   string `gorm:"text;not null;" json:"role_label"`
	AccessLevel int    `gorm:"text;not null;" json:"access_level"`
}

type UserRoles interface {
	getRole() UserRole
}

func GetUserRoleData() []UserRole {
	return []UserRole{
		{RoleLabel: "Admin", AccessLevel: 1},
		{RoleLabel: "Moderator", AccessLevel: 2},
		{RoleLabel: "User", AccessLevel: 3},
	}
}
