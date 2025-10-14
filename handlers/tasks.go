package handlers

import (
	"html/template"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"task-manager/database"
	"task-manager/models"
)

// =======================
// LIST TASKS (KANBAN)
// =======================
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Preload("AssignedUser").Preload("Creator").Find(&tasks)

	sess, _ := store.Get(r, "session-name")
	successFlashes := sess.Flashes("success")
	errorFlashes := sess.Flashes("error")
	sess.Save(r, w)

	// Ambil URL router
	var dashboardURL, tasksURL, usersURL, createTaskURL string
	if route := Router.Get("dashboard"); route != nil {
		u, _ := route.URL()
		dashboardURL = u.String()
	}
	if route := Router.Get("tasks"); route != nil {
		u, _ := route.URL()
		tasksURL = u.String()
	}
	if route := Router.Get("users"); route != nil {
		u, _ := route.URL()
		usersURL = u.String()
	}
	if route := Router.Get("tasks.create"); route != nil {
		u, _ := route.URL()
		createTaskURL = u.String()
	}

	data := map[string]interface{}{
		"Tasks":        tasks,
		"CreateTaskURL":  createTaskURL,
		"DashboardURL":   dashboardURL,
		"TasksURL":  tasksURL,
		"UsersURL":       usersURL,
		"CurrentPath":  r.URL.Path,
		"Success":      successFlashes,
		"Error":        errorFlashes,
	}

	funcs := template.FuncMap{
		"year": func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },

		// bikin list manual karena Go template nggak punya bawaan
		"list": func(values ...string) []string {
			return values
		},
		"title": func(s string) string {
			return strings.Title(s)
		},
		"filterStatus": func(tasks []models.Task, status string) []models.Task {
			var result []models.Task
			for _, t := range tasks {
				if strings.ToLower(t.Status) == strings.ToLower(status) {
					result = append(result, t)
				}
			}
			return result
		},
		"hasTask": func(tasks []models.Task, status string) bool {
			for _, t := range tasks {
				if strings.ToLower(t.Status) == strings.ToLower(status) {
					return true
				}
			}
			return false
		},
		"substr": func(s string, start, end int) string {
			if len(s) == 0 {
				return "?"
			}
			if start < 0 {
				start = 0
			}
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"statusColor": func(status string) string {
			switch strings.ToLower(status) {
			case "ready":
				return "border-blue-400"
			case "in progress":
				return "border-yellow-400"
			case "fix":
				return "border-red-400"
			case "code review":
				return "border-indigo-400"
			case "test":
				return "border-orange-400"
			case "done":
				return "border-green-400"
			default:
				return "border-gray-300"
			}
		},
		"typeColor": func(t string) string {
			switch strings.ToLower(t) {
			case "backend":
				return "bg-purple-100 text-purple-700"
			case "frontend":
				return "bg-blue-100 text-blue-700"
			case "qa":
				return "bg-pink-100 text-pink-700"
			default:
				return "bg-gray-100 text-gray-700"
			}
		},
		"priorityColor": func(p string) string {
			switch strings.ToLower(p) {
			case "high":
				return "text-red-500"
			case "medium":
				return "text-yellow-500"
			case "low":
				return "text-green-500"
			default:
				return "text-gray-500"
			}
		},
	}

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

// ========================
// SHOW CREATE FORM
// ========================
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	database.DB.Find(&users)

	var dashboardURL, usersURL, tasksURL, storeTaskURL string
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
	if route := Router.Get("tasks.store"); route != nil {
		u, _ := route.URL()
		storeTaskURL = u.String()
	}

	data := map[string]interface{}{
		"Users":         users,
		"DashboardURL":  dashboardURL,
		"TasksURL":      tasksURL,
		"UsersURL":       usersURL,
		"StoreTaskURL":  storeTaskURL,
		"CurrentPath":   r.URL.Path,
	}

	funcs := template.FuncMap{
		"year":      func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
	}

	tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
		"templates/layouts/base.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
		"templates/tasks/create.html",
	))

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ========================
// STORE TASK (POST)
// ========================
func StoreTaskHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "session-name")

	if err := r.ParseForm(); err != nil {
		sess.AddFlash("Form tidak valid: "+err.Error(), "error")
		sess.Save(r, w)
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	// Parse tanggal
	layout := "2006-01-02"
	var startDate, dueDate *time.Time
	if s := strings.TrimSpace(r.FormValue("start_date")); s != "" {
		if t, err := time.Parse(layout, s); err == nil {
			startDate = &t
		}
	}
	if s := strings.TrimSpace(r.FormValue("due_date")); s != "" {
		if t, err := time.Parse(layout, s); err == nil {
			dueDate = &t
		}
	}

	// Parse AssignedTo
	var assignedTo uint
	if at := strings.TrimSpace(r.FormValue("assigned_to")); at != "" {
		if parsed, err := strconv.ParseUint(at, 10, 32); err == nil {
			assignedTo = uint(parsed)
		}
	}

	// Ambil CreatedBy dari session user_id
	var createdBy uint = 1
	if uid, ok := sess.Values["user_id"].(uint); ok {
		createdBy = uid
	} else if uidStr, ok := sess.Values["user_id"].(string); ok {
		if parsed, err := strconv.ParseUint(uidStr, 10, 32); err == nil {
			createdBy = uint(parsed)
		}
	}

	task := models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Type:        r.FormValue("type"),
		Priority:    r.FormValue("priority"),
		Status:      r.FormValue("status"),
		TaskLink:    r.FormValue("task_link"),
		StartDate:   startDate,
		DueDate:     dueDate,
		AssignedTo:  assignedTo,
		CreatedBy:   createdBy,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		sess.AddFlash("Gagal menyimpan task: "+err.Error(), "error")
		sess.Save(r, w)
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	sess.AddFlash("Task berhasil dibuat", "success")
	sess.Save(r, w)
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
