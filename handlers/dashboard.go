package handlers

import (
	"html/template"
	"net/http"
	"task-manager/database"
	"task-manager/models"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

var Router *mux.Router

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Preload("AssignedUser").Preload("Creator").Find(&tasks)

	// Hitung task berdasarkan status
	var completedCount, inProgressCount int
	for _, t := range tasks {
		switch t.Status {
		case "done":
			completedCount++
		case "in progress":
			inProgressCount++
		}
	}

	// Ambil total user dari database
	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	// Ambil URL dari named route
	var dashboardURL, usersURL string

	if route := Router.Get("dashboard"); route != nil {
		u, _ := route.URL()
		dashboardURL = u.String()
	}
	if route := Router.Get("users"); route != nil {
		u, _ := route.URL()
		usersURL = u.String()
	}

	// Kirim semua data ke template
	data := map[string]interface{}{
		"Tasks":         tasks,
		"TotalTasks":    len(tasks),
		"Completed":     completedCount,
		"InProgress":    inProgressCount,
		"TotalUsers":    totalUsers,
		"CurrentPath":   r.URL.Path,
		"DashboardURL":  dashboardURL,
		"UsersURL":      usersURL,
	}


	tmpl := template.New("").Funcs(template.FuncMap{
		"year":      func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
	})

	tmpl = template.Must(tmpl.ParseFiles(
		"templates/dashboard.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
	))

	err := tmpl.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
