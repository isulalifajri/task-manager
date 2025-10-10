package models

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100)"`
	Role string `gorm:"type:varchar(50)"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(100)"`
	Description string
	Status      string `gorm:"type:varchar(20)"`
	AssignedTo  uint
}
