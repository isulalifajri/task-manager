package database

import (
	"fmt"
	"time"
	"task-manager/models"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	// ===== Roles =====
	roles := []models.Role{
		{Name: "superadmin", Description: "Akses penuh"},
		{Name: "manager", Description: "Project Manager"},
		{Name: "developer", Description: "Developer"},
		{Name: "reviewer", Description: "Code Reviewer"},
	}
	DB.Create(&roles)
	fmt.Println("Roles siap!")

	// ===== Users =====
	users := []models.User{
		{Name: "Hening", Username: "hening", Email: "hening@example.com", Password: hashPassword("123456"), RoleID: 2},
		{Name: "Dwi", Username: "dwi", Email: "dwi@example.com", Password: hashPassword("123456"), RoleID: 3},
		{Name: "Raka", Username: "raka", Email: "raka@example.com", Password: hashPassword("123456"), RoleID: 4},
		{Name: "Langit", Username: "langit", Email: "langit@example.com", Password: hashPassword("123456"), RoleID: 1},
		{Name: "Flower", Username: "flower", Email: "flower@example.com", Password: hashPassword("123456"), RoleID: 3},
		{Name: "Melati", Username: "melati", Email: "melati@example.com", Password: hashPassword("123456"), RoleID: 4},
	}
	DB.Create(&users)
	fmt.Println("Users siap!")

	// ===== Tasks =====
	now := time.Now()
	nextWeek := now.AddDate(0, 0, 7)

	tasks := []models.Task{
		{
			Title:       "Setup project structure",
			Description: "Inisialisasi repo dan struktur folder",
			Status:      "ready",
			Type:        "Backend",
			Priority:    "High",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  2,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/1",
		},
		{
			Title:       "Implementasi CRUD Task",
			Description: "Buat API untuk create, update, delete task",
			Status:      "in progress",
			Type:        "Backend",
			Priority:    "Critical",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  2,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/2",
		},
		{
			Title:       "Perbaiki validasi input",
			Description: "Bug pada form task input",
			Status:      "fix",
			Type:        "Frontend",
			Priority:    "Medium",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  3,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/3",
		},
		{
			Title:       "Review pull request API",
			Description: "Periksa PR untuk modul task",
			Status:      "code review",
			Type:        "Backend",
			Priority:    "High",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  4,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/4",
		},
		{
			Title:       "Testing fitur login",
			Description: "QA untuk fitur login dan autentikasi",
			Status:      "test",
			Type:        "QA",
			Priority:    "Medium",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  5,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/5",
		},
		{
			Title:       "Deploy ke staging",
			Description: "Pastikan semua task sudah siap sebelum deploy",
			Status:      "done",
			Type:        "DevOps",
			Priority:    "Low",
			StartDate:   &now,
			DueDate:     &nextWeek,
			AssignedTo:  6,
			CreatedBy:   2,
			TaskLink:    "https://example.com/task/6",
		},
	}
	DB.Create(&tasks)
	fmt.Println("Tasks siap!")
}

func hashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		panic("Gagal hash password")
	}
	return string(hash)
}
