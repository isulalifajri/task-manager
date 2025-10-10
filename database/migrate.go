package database

import (
	"fmt"
	"task-manager/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Task{},
	)
	if err != nil {
		fmt.Println("Gagal migrasi:", err)
		return
	}
	fmt.Println("Migrasi berhasil!")
}
