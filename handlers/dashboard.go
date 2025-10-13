package handlers

import (
	"html/template"
	"net/http"
	"task-manager/database"
	"task-manager/models"
	"runtime"
    "strings"
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

    // Hitung total user
    var totalUsers int64
    database.DB.Model(&models.User{}).Count(&totalUsers)

    // Ambil URL dari router global
    var dashboardURL, usersURL, tasksURL string
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

    data := map[string]interface{}{
        "Tasks":        tasks,
        "TotalTasks":   len(tasks),
        "Completed":    completedCount,
        "InProgress":   inProgressCount,
        "TotalUsers":   totalUsers,
        "CurrentPath":  r.URL.Path,
        "DashboardURL": dashboardURL,
        "UsersURL":     usersURL,
        "TasksURL":  tasksURL,
    }

    funcs := template.FuncMap{
        "year":      func() int { return time.Now().Year() },
        "goversion": func() string { return runtime.Version() },
        "hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
    }

    tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
        "templates/layouts/base.html",
        "templates/layouts/header.html",
        "templates/layouts/sidebar.html",
        "templates/layouts/footer.html",
        "templates/dashboard.html",
    ))

    // Gunakan base layout
    err := tmpl.ExecuteTemplate(w, "base", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

