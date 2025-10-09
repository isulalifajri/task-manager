package models

// User merepresentasikan data pengguna sistem
type User struct {
	ID    int
	Name  string
	Role  string // contoh: "manager", "developer", "reviewer"
}

// Task merepresentasikan tugas di sistem
type Task struct {
	ID          int
	Title       string
	Description string
	Status      string // contoh: "ready", "on progress", "review", "done"
	AssignedTo  int    // ID user yang ditugaskan
}
