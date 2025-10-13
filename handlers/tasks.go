package handlers

import (
	"html/template"
	"net/http"
	"strings"
	"time"
	"task-manager/database"
	"task-manager/models"
	"runtime"
)

// TasksHandler menampilkan semua task
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil semua task
	var tasks []models.Task
	database.DB.Find(&tasks)

	// Ambil URL dari router untuk sidebar
	var dashboardURL, tasksURL, usersURL string
	if route := Router.Get("dashboard"); route != nil {
		u, _ := route.URL()
		dashboardURL = u.String()
	}
	if route := Router.Get("users"); route != nil {
        u, _ := route.URL()
        usersURL = u.String()
    }
	if route := Router.Get("tasks"); route != nil {
		u, _ := route.URL()
		tasksURL = u.String()
	}

	// Data untuk template
	data := map[string]interface{}{
		"Tasks":        tasks,
		"CurrentPath":  r.URL.Path,
		"DashboardURL": dashboardURL,
		"UsersURL":     usersURL,
		"TasksURL":     tasksURL,
	}

	// Template functions
	funcs := template.FuncMap{
		"year":      func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
	}

	// Load template
	tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
		"templates/layouts/base.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
		"templates/tasks/tasks.html",
	))

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
