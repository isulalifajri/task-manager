package database

import (
	"fmt"
	"task-manager/models"

	"golang.org/x/crypto/bcrypt"
)

// Seed awal role, user, dan task
func Seed() {
	// ===== Roles =====
	roles := []models.Role{
		{Name: "superadmin", Description: "Akses penuh"},
		{Name: "manager", Description: "Project Manager"},
		{Name: "developer", Description: "Developer"},
		{Name: "reviewer", Description: "Code Reviewer"},
	}

	for _, role := range roles {
		var existing models.Role
		if err := DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			DB.Create(&role)
		}
	}

	fmt.Println("Roles siap!")

	// ===== Users =====
	users := []models.User{
		{Name: "Hening", Username: "hening", Email: "hening@example.com", Password: hashPassword("123456"), RoleID: 2},
		{Name: "Dwi", Username: "dwi", Email: "dwi@example.com", Password: hashPassword("123456"), RoleID: 3},
		{Name: "Raka", Username: "raka", Email: "raka@example.com", Password: hashPassword("123456"), RoleID: 4},
	}

	for _, user := range users {
		var existing models.User
		if err := DB.Where("email = ?", user.Email).First(&existing).Error; err != nil {
			DB.Create(&user)
		}
	}

	fmt.Println("Users siap!")

	// ===== Tasks =====
	tasks := []models.Task{
		{
			Title:       "Setup project",
			Description: "Inisialisasi repo dan struktur folder",
			Status:      "done",
			AssignedTo:  2, // ID user
			CreatedBy:   2, // ID user yang membuat task
			TaskLink:    "https://example.com/task/1",
		},
		{
			Title:       "Buat API task",
			Description: "Implementasi CRUD task",
			Status:      "in progress",
			AssignedTo:  2,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/2",
		},
		{
			Title:       "Review code",
			Description: "Periksa pull request",
			Status:      "ready",
			AssignedTo:  3,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/3",
		},
	}

	for _, task := range tasks {
		var existing models.Task
		if err := DB.Where("title = ?", task.Title).First(&existing).Error; err != nil {
			DB.Create(&task)
		}
	}

	fmt.Println("Tasks siap!")
}

// hashPassword bikin password terenkripsi
func hashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic("Gagal hash password")
	}
	return string(hash)
}
