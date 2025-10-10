package handlers

import (
	"html/template"
	"net/http"
	"task-manager/database"
	"task-manager/models"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Preload("AssignedUser").Preload("Creator").Find(&tasks)

	tmpl := template.Must(template.ParseFiles(
		"templates/dashboard.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
	))

	err := tmpl.ExecuteTemplate(w, "dashboard.html", tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
