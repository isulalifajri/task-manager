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
	database.DB.
		Preload("AssignedUser").
		Preload("Creator").
		Find(&tasks)


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
		"year": func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },

		// bikin list manual karena Go template nggak punya bawaan
		"list": func(values ...string) []string {
			return values
		},

		// biar bisa filter task per status
		"filterStatus": func(tasks []models.Task, status string) []models.Task {
			var filtered []models.Task
			for _, t := range tasks {
				if strings.EqualFold(t.Status, status) {
					filtered = append(filtered, t)
				}
			}
			return filtered
		},

		// cek apakah ada task di status tertentu
		"hasTask": func(tasks []models.Task, status string) bool {
			for _, t := range tasks {
				if strings.EqualFold(t.Status, status) {
					return true
				}
			}
			return false
		},

		// biar capitalize
		"title": strings.Title,

		// ambil inisial
		"substr": func(s string, start, length int) string {
			if len(s) < start {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},

		// warna border berdasarkan status
		"statusColor": func(status string) string {
			switch strings.ToLower(status) {
			case "ready":
				return "border-gray-300"
			case "in progress":
				return "border-yellow-400"
			case "fix":
				return "border-red-400"
			case "code review":
				return "border-orange-400"
			case "test":
				return "border-blue-400"
			case "done":
				return "border-green-400"
			default:
				return "border-gray-200"
			}
		},

		// warna type label
		"typeColor": func(tp string) string {
			switch strings.ToLower(tp) {
			case "backend":
				return "bg-blue-100 text-blue-700"
			case "frontend":
				return "bg-pink-100 text-pink-700"
			case "qa":
				return "bg-yellow-100 text-yellow-700"
			case "devops":
				return "bg-green-100 text-green-700"
			default:
				return "bg-gray-100 text-gray-700"
			}
		},

		// warna priority
		"priorityColor": func(priority string) string {
			switch strings.ToLower(priority) {
			case "low":
				return "text-green-600"
			case "medium":
				return "text-yellow-600"
			case "high":
				return "text-orange-600"
			case "critical":
				return "text-red-600 font-semibold"
			default:
				return "text-gray-500"
			}
		},
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
