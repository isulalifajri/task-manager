package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'ready'"`
	TaskLink    string
	AssignedTo  uint
	AssignedUser User `gorm:"foreignKey:AssignedTo"` // relasi ke User
	CreatedBy   uint
	Creator      User `gorm:"foreignKey:CreatedBy"`  // relasi ke User
}

