package handlers

import (
    "fmt"
    "net/http"
)

// HomeHandler — untuk halaman utama
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Selamat datang di Task Manager API 🎯")
}

// TaskHandler — untuk endpoint /tasks
func TaskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        fmt.Fprintln(w, "Menampilkan semua task")
    case "POST":
        fmt.Fprintln(w, "Menambahkan task baru")
    default:
        http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
    }
}
