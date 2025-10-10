package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	RoleID   uint
	Role     Role
	Tasks    []Task `gorm:"foreignKey:CreatedBy"`
}
