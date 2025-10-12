package handlers

import (
	"html/template"
	"net/http"
	"task-manager/database"
	"task-manager/models"
	"runtime"
	"time"
	"strconv"
	"math"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
    // Ambil query params untuk pagination
    pageStr := r.URL.Query().Get("page")
    const pageSize = 5 // tampilkan 5 per halaman
    page := 1
    if pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }

    // Hitung offset
    offset := (page - 1) * pageSize

    // Ambil total user (buat pagination info)
    var total int64
    database.DB.Model(&models.User{}).Count(&total)

    // Ambil user dengan limit + offset
    var users []models.User
    database.DB.Preload("Role").
        Limit(pageSize).
        Offset(offset).
        Find(&users)

    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

    // Ambil URL dashboard (biar sidebar tetap pakai)
    var dashboardURL, usersURL string
    if route := Router.Get("dashboard"); route != nil {
        u, _ := route.URL()
        dashboardURL = u.String()
    }
    if route := Router.Get("users"); route != nil {
        u, _ := route.URL()
        usersURL = u.String()
    }

    data := map[string]interface{}{
        "Users":        users,
        "CurrentPath":  r.URL.Path,
        "DashboardURL": dashboardURL,
        "UsersURL":     usersURL,
        "Page":         page,
        "TotalPages":   totalPages,
    }

    tmpl := template.New("").Funcs(template.FuncMap{
        "year":      func() int { return time.Now().Year() },
        "goversion": func() string { return runtime.Version() },
        "add":       func(a, b int) int { return a + b },
        "sub":       func(a, b int) int { return a - b },
		"until":     func(n int) []int {
			a := make([]int, n)
			for i := 0; i < n; i++ {
				a[i] = i
			}
			return a
		},
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
