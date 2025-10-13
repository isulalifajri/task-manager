package models

import (
	"time"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title        string     `gorm:"not null"`
	Description  string
	Status       string     `gorm:"default:'ready'"` // ready, in progress, fix, code review, test, done
	Type         string     `gorm:"size:50"`         // Backend, Frontend, QA, DevOps, dll
	Priority     string     `gorm:"size:20"`         // Low, Medium, High, Critical
	TaskLink     string
	StartDate    *time.Time
	DueDate      *time.Time
	AssignedTo   uint
	AssignedUser User `gorm:"foreignKey:AssignedTo"`
	CreatedBy    uint
	Creator      User `gorm:"foreignKey:CreatedBy"`
}
