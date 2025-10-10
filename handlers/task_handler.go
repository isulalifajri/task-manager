package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"task-manager/database"
	"task-manager/models"
)

// ==================== API JSON ====================

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Selamat datang di Task Manager API 🎯"))
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		var tasks []models.Task
		if err := database.DB.Preload("AssignedUser").Preload("Creator").Find(&tasks).Error; err != nil {
			http.Error(w, "Gagal mengambil task", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tasks)

	case "POST":
		var newTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		if err := database.DB.Create(&newTask).Error; err != nil {
			http.Error(w, "Gagal membuat task", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(newTask)

	case "PUT":
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		var task models.Task
		if err := database.DB.First(&task, id).Error; err != nil {
			http.Error(w, "Task tidak ditemukan", http.StatusNotFound)
			return
		}

		var updatedTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		task.Title = updatedTask.Title
		task.Description = updatedTask.Description
		task.Status = updatedTask.Status
		task.TaskLink = updatedTask.TaskLink
		task.AssignedTo = updatedTask.AssignedTo
		task.CreatedBy = updatedTask.CreatedBy

		if err := database.DB.Save(&task).Error; err != nil {
			http.Error(w, "Gagal update task", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(task)

	case "DELETE":
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		if err := database.DB.Delete(&models.Task{}, id).Error; err != nil {
			http.Error(w, "Gagal hapus task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Task dihapus"})

	default:
		http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}

// ==================== HTML UI ====================

func TaskHTMLHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Preload("AssignedUser").Preload("Creator").Find(&tasks)

	tmpl := template.Must(template.ParseFiles("templates/layout/header.html", "templates/layout/footer.html", "templates/tasks.html"))
	tmpl.ExecuteTemplate(w, "base", tasks)
}

func TaskFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout/header.html", "templates/layout/footer.html", "templates/task_form.html"))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	status := r.FormValue("status")
	assignedTo, _ := strconv.Atoi(r.FormValue("assigned_to"))
	createdBy, _ := strconv.Atoi(r.FormValue("created_by"))

	task := models.Task{
		Title:       title,
		Description: description,
		Status:      status,
		AssignedTo:  uint(assignedTo),
		CreatedBy:   uint(createdBy),
	}

	if err := database.DB.Create(&task).Error; err != nil {
		http.Error(w, "Gagal membuat task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks/html", http.StatusSeeOther)
}

// ✅ Handler tambahan: DELETE via tombol HTML
func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/tasks/delete/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := database.DB.Delete(&models.Task{}, id).Error; err != nil {
		http.Error(w, "Gagal menghapus task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks/html", http.StatusSeeOther)
}