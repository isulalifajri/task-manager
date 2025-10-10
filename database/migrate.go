package database

import (
	"fmt"
	"task-manager/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		panic("Gagal migrasi database: " + err.Error())
	}
	fmt.Println("Migrasi database selesai (GORM)")
}
