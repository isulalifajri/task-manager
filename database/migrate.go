package database

import (
	"fmt"
	"task-manager/models"
)

// Jalankan migrasi
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

// Tambahkan fungsi baru ini:
func Fresh() {
	fmt.Println("Menghapus semua tabel lama...")
	DB.Migrator().DropTable(&models.Task{}, &models.User{}, &models.Role{})

	fmt.Println("Migrasi ulang dan seeding data...")
	Migrate()
	Seed()
	fmt.Println("Database sudah direfresh dan diisi ulang!")
}
