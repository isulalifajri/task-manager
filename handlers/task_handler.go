package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager/models"
)

// data sementara
var users = []models.User{
	{ID: 1, Name: "Hening", Role: "manager"},
	{ID: 2, Name: "Dwi", Role: "developer"},
	{ID: 3, Name: "Raka", Role: "reviewer"},
}

var tasks = []models.Task{
	{ID: 1, Title: "Setup project", Description: "Inisialisasi repo dan struktur folder", Status: "done", AssignedTo: 2},
	{ID: 2, Title: "Buat API task", Description: "Implementasi CRUD task", Status: "on progress", AssignedTo: 2},
	{ID: 3, Title: "Review code", Description: "Periksa pull request", Status: "ready", AssignedTo: 3},
}

// HomeHandler — untuk halaman utama
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Selamat datang di Task Manager API"))
}

// TaskHandler — endpoint /tasks
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			json.NewEncoder(w).Encode(tasks)
			return
		}

		intID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}
		id := uint(intID)

		for _, task := range tasks {
			if task.ID == id {
				json.NewEncoder(w).Encode(task)
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	case "POST":
		var newTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		newTask.ID = uint(len(tasks) + 1)
		tasks = append(tasks, newTask)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	case "PUT":
		idParam := r.URL.Query().Get("id")
		intID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}
		id := uint(intID)

		var updatedTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		for i, t := range tasks {
			if t.ID == id {
				updatedTask.ID = id
				tasks[i] = updatedTask
				json.NewEncoder(w).Encode(updatedTask)
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	case "DELETE":
		idParam := r.URL.Query().Get("id")
		intID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}
		id := uint(intID)

		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Task dihapus"})
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	default:
		http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}
