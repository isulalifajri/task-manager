package handlers

import (
	"html/template"
	"net/http"
	"task-manager/database"
	"task-manager/models"

	"github.com/gorilla/mux"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil semua user dari database
	var users []models.User
	database.DB.Preload("Role").Find(&users)

	// Ambil URL dashboard (biar sidebar tetap bisa pakai)
	var dashboardURL string
	if route := Router.Get("dashboard"); route != nil {
		u, _ := route.URL()
		dashboardURL = u.String()
	}

	// Ambil URL users (biar sidebar pakai {{.UsersURL}})
	var usersURL string
	if route := Router.Get("users"); route != nil {
		u, _ := route.URL()
		usersURL = u.String()
	}

	// Data untuk dikirim ke template
	data := map[string]interface{}{
		"Users":        users,
		"CurrentPath":  r.URL.Path,
		"DashboardURL": dashboardURL,
		"UsersURL":     usersURL,
	}

	// Load template
	tmpl := template.New("").Funcs(template.FuncMap{
		"year":      func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
	})

	tmpl = template.Must(tmpl.ParseFiles(
		"templates/users.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
	))

	err := tmpl.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
