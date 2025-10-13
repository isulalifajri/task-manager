package handlers

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	"task-manager/database"
	"task-manager/models"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
    "github.com/gorilla/sessions"

)

// Session store global
var store = sessions.NewCookieStore([]byte("super-secret-key")) // Ganti dengan key aman

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil query params untuk pagination
	pageStr := r.URL.Query().Get("page")
	const limit = 5
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	offset := (page - 1) * limit

	// Ambil total user
	var total int64
	database.DB.Model(&models.User{}).Count(&total)

	// Ambil data user dengan relasi role
	var users []models.User
	database.DB.Preload("Role").Limit(limit).Offset(offset).Find(&users)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Ambil URL dari router
	var dashboardURL, usersURL, createUsersURL, editUsersURL, deleteUserURL string
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
	if route := Router.Get("users.edit"); route != nil {
		u, _ := route.URL("id", "0")
		editUsersURL = u.String()
	}
	if route := Router.Get("users.delete"); route != nil {
		u, _ := route.URL("id", "0")
		deleteUserURL = u.String()
	}

	// Ambil flash messages
	sess, _ := store.Get(r, "session-name")
	successFlashes := sess.Flashes("success")
	errorFlashes := sess.Flashes("error")
	sess.Save(r, w)

	// Data untuk template
	data := map[string]interface{}{
		"Users":          users,
		"CurrentPath":    r.URL.Path,
		"DashboardURL":   dashboardURL,
		"UsersURL":       usersURL,
		"CreateUsersURL": createUsersURL,
		"EditUserURL":    editUsersURL,
		"DeleteUserURL":  deleteUserURL,
		"Page":           page,
		"Limit":          limit,
		"TotalPages":     totalPages,
		"Offset":         offset,
		"Success":        successFlashes,
		"Error":          errorFlashes,
	}

	// Template functions
	funcs := template.FuncMap{
		"year": func() int { return time.Now().Year() },
		"goversion": func() string { return runtime.Version() },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"until": func(n int) []int {
			a := make([]int, n)
			for i := 0; i < n; i++ {
				a[i] = i
			}
			return a
		},
		"editURL": func(base string, id uint) string {
			return strings.Replace(base, "0", fmt.Sprintf("%d", id), 1)
		},
		"deleteURL": func(base string, id uint) string {
			return strings.Replace(base, "0", fmt.Sprintf("%d", id), 1)
		},
		"hasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
	}

	// Load template
	tmpl := template.Must(template.New("base.html").Funcs(funcs).ParseFiles(
		"templates/layouts/base.html",
		"templates/layouts/header.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/footer.html",
		"templates/users/users.html",
	))

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Ambil URL dashboard & users (untuk sidebar)
    var dashboardURL, usersURL, storeUsersURL string
    if route := Router.Get("dashboard"); route != nil {
        u, _ := route.URL()
        dashboardURL = u.String()
    }
    if route := Router.Get("users"); route != nil {
        u, _ := route.URL()
        usersURL = u.String()
    }
	if route := Router.Get("users.store"); route != nil {
        u, _ := route.URL()
        storeUsersURL = u.String()
    }

	var roles []models.Role
    if err := database.DB.Find(&roles).Error; err != nil {
        http.Error(w, "Failed to get roles", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "CurrentPath":  r.URL.Path,
        "DashboardURL": dashboardURL,
        "UsersURL":     usersURL,
        "StoreUsersURL": storeUsersURL,
		"Roles":        roles,
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
    sess, _ := store.Get(r, "session-name") // ambil session

    if err := r.ParseForm(); err != nil {
        sess.AddFlash("Form tidak valid: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    username := r.FormValue("username")

    // Cek dulu username unik
    var existingUser models.User
    if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
        sess.AddFlash("Username sudah digunakan, coba yang lain.", "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    } else if err != gorm.ErrRecordNotFound {
        sess.AddFlash("Gagal mengecek username: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    password := r.FormValue("password")
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        sess.AddFlash("Gagal memproses password: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    roleID, _ := strconv.Atoi(r.FormValue("role_id"))

    user := models.User{
        Name:     r.FormValue("name"),
        Username: username,
        Email:    r.FormValue("email"),
        Password: string(hashedPassword),
        RoleID:   uint(roleID),
    }

    if err := database.DB.Create(&user).Error; err != nil {
        sess.AddFlash("Gagal menyimpan user: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    sess.AddFlash("User berhasil dibuat", "success")
    sess.Save(r, w)
    http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Ambil URL dashboard, users, dan store (update)
    var dashboardURL, usersURL, updateUserURL string
    if route := Router.Get("dashboard"); route != nil {
        u, _ := route.URL()
        dashboardURL = u.String()
    }
    if route := Router.Get("users"); route != nil {
        u, _ := route.URL()
        usersURL = u.String()
    }
    if route := Router.Get("users.update"); route != nil {
        u, _ := route.URL("id", idStr)
        updateUserURL = u.String()
    }

    // Ambil semua roles dari DB
    var roles []models.Role
    if err := database.DB.Find(&roles).Error; err != nil {
        http.Error(w, "Failed to get roles", http.StatusInternalServerError)
        return
    }

    // Ambil user berdasarkan ID
    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Data untuk template
    data := map[string]interface{}{
        "CurrentPath":   r.URL.Path,
        "DashboardURL":  dashboardURL,
        "UsersURL":      usersURL,
        "UpdateUserURL": updateUserURL,
        "Roles":         roles,
        "User":          user,
        "Errors": map[string]string{},
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
        "templates/users/user_form.html",
    ))

    if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    sess, _ := store.Get(r, "session-name")

    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        sess.AddFlash("ID user tidak valid: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    if err := r.ParseForm(); err != nil {
        sess.AddFlash("Form tidak valid: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        sess.AddFlash("User tidak ditemukan: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    user.Name = r.FormValue("name")
    user.Username = r.FormValue("username")
    user.Email = r.FormValue("email")

    roleID, _ := strconv.Atoi(r.FormValue("role_id"))
    user.RoleID = uint(roleID)

    password := r.FormValue("password")
    if password != "" {
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            sess.AddFlash("Gagal mengubah password: "+err.Error(), "error")
            sess.Save(r, w)
            http.Redirect(w, r, "/users", http.StatusSeeOther)
            return
        }
        user.Password = string(hashed)
    }

    if err := database.DB.Save(&user).Error; err != nil {
        sess.AddFlash("Gagal memperbarui user: "+err.Error(), "error")
        sess.Save(r, w)
        http.Redirect(w, r, "/users", http.StatusSeeOther)
        return
    }

    sess.AddFlash("User berhasil diperbarui", "success")
    sess.Save(r, w)
    http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	sess, _ := store.Get(r, "session-name")

	if err != nil {
		sess.AddFlash("ID user tidak valid", "error")
		sess.Save(r, w)
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		sess.AddFlash("Gagal menghapus user", "error")
		sess.Save(r, w)
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	sess.AddFlash("User berhasil dihapus", "success")
	sess.Save(r, w)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}



