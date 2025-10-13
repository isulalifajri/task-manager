package handlers

import (
	"html/template"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"task-manager/database"
	"task-manager/models"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil query params untuk pagination
	pageStr := r.URL.Query().Get("page")
	const limit = 5 // jumlah data per halaman
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit

	// Hitung total data
	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	// Ambil data user (dengan relasi role)
	var users []models.User
	database.DB.Preload("Role").
		Limit(limit).
		Offset(offset).
		Find(&users)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Ambil URL dashboard & users dari router (untuk sidebar)
	var dashboardURL, usersURL, createUsersURL string
	if route := Router.Get("dashboard"); route != nil {
		u, _ := route.URL()
		dashboardURL = u.String()
	}
	if route := Router.Get("users"); route != nil {
		u, _ := route.URL()
		usersURL = u.String()
	}
	if route := Router.Get("users.create"); route != nil {
		u, _ := route.URL()
		createUsersURL = u.String()
	}

	// Data untuk template
	data := map[string]interface{}{
		"Users":        users,
		"CurrentPath":  r.URL.Path,
		"DashboardURL": dashboardURL,
		"UsersURL":     usersURL,
		"CreateUsersURL": createUsersURL,
		"Page":         page,
		"Limit":        limit,
		"TotalPages":   totalPages,
		"Offset":       offset,
	}

	// Template functions
	funcs := template.FuncMap{
		"year":      func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a - b },
		"mul":       func(a, b int) int { return a * b },
		"until": func(n int) []int {
			a := make([]int, n)
			for i := 0; i < n; i++ {
				a[i] = i
			}
			return a
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	}

	// Load semua template (base layout + komponen)
	tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
		"templates/layouts/base.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
		"templates/users/users.html",
	))

	// Gunakan base layout
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Ambil URL dashboard & users (untuk sidebar)
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
        "CurrentPath":  r.URL.Path,
        "DashboardURL": dashboardURL,
        "UsersURL":     usersURL,
    }

    funcs := template.FuncMap{
        "year":      func() int { return time.Now().Year() },
        "goversion": func() string { return runtime.Version() },
        "add":       func(a, b int) int { return a + b },
        "hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
    }

    tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
        "templates/layouts/base.html",
        "templates/layouts/header.html",
        "templates/layouts/sidebar.html",
        "templates/layouts/footer.html",
        "templates/users/user_create.html",
    ))

    err := tmpl.ExecuteTemplate(w, "base", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func StoreUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     r.FormValue("name"),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
	}

	// Ambil role dari form
	roleID, _ := strconv.Atoi(r.FormValue("role_id"))
	user.RoleID = uint(roleID)

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

