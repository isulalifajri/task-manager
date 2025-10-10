## Install Golang

Download golang : `https://go.dev/dl/`

cek instalasi : `go version`

## inisialisi Modul GO
Setiap project Go butuh modul (semacam package.json kalau di Node.js).

jalankan ini:

```
go mod init task-manager

```

setelah itu nanti akan muncul file: `go.mod`

# Start Projects

Buat file dengan nama 'main.go' dg isi code seperti ini:

```
package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager/handlers"
	"task-manager/middlewares"
)

func main() {
	mux := http.NewServeMux()

	// route
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/tasks", handlers.TaskHandler)

	// pasang middleware (urutan penting)
	handlerWithMiddleware := middlewares.LoggingMiddleware(middlewares.CORSMiddleware(mux))

	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}

```

Kemudian buat folder dengan nama `handlers` dan buat file di dalamnya dengan nama `task_handlers.go` dan isi filenya seperti ini:

```
package handlers

import (
	"encoding/json"
	"net/http"
	"task-manager/models"
)

// data sementara (dummy)
var users = []models.User{
	{ID: 1, Name: "Hening", Role: "manager"},
	{ID: 2, Name: "Dwi", Role: "developer"},
	{ID: 3, Name: "Raka", Role: "reviewer"},
}

var tasks = []models.Task{
	{ID: 1, Title: "Setup project", Description: "Inisialisasi repo dan struktur folder", Status: "done", AssignedTo: 2},
	{ID: 2, Title: "Buat API task", Description: "Implementasi CRUD task", Status: "on progress", AssignedTo: 2},
	{ID: 3, Title: "Review code", Description: "Periksa pull request", Status: "ready", AssignedTo: 3},
}

// HomeHandler — untuk halaman utama
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Selamat datang di Task Manager API"))
}

// TaskHandler — endpoint /tasks
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(tasks)
	case "POST":
		http.Error(w, "Fitur tambah task belum dibuat", http.StatusNotImplemented)
	default:
		http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}


```

Buat folder models dan didalamnya buat file dengan nama `models.go` isi filenya dengan code ini:

```
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager/models"
)

// data sementara
var users = []models.User{
	{ID: 1, Name: "Hening", Role: "manager"},
	{ID: 2, Name: "Dwi", Role: "developer"},
	{ID: 3, Name: "Raka", Role: "reviewer"},
}

var tasks = []models.Task{
	{ID: 1, Title: "Setup project", Description: "Inisialisasi repo dan struktur folder", Status: "done", AssignedTo: 2},
	{ID: 2, Title: "Buat API task", Description: "Implementasi CRUD task", Status: "on progress", AssignedTo: 2},
	{ID: 3, Title: "Review code", Description: "Periksa pull request", Status: "ready", AssignedTo: 3},
}

// HomeHandler — untuk halaman utama
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Selamat datang di Task Manager API"))
}

// TaskHandler — endpoint /tasks
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			// tanpa parameter: tampilkan semua
			json.NewEncoder(w).Encode(tasks)
			return
		}

		// kalau ada parameter id
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		for _, task := range tasks {
			if task.ID == id {
				json.NewEncoder(w).Encode(task)
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	case "POST":
		var newTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		newTask.ID = len(tasks) + 1
		tasks = append(tasks, newTask)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	case "PUT":
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		var updatedTask models.Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		for i, t := range tasks {
			if t.ID == id {
				updatedTask.ID = id
				tasks[i] = updatedTask
				json.NewEncoder(w).Encode(updatedTask)
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	case "DELETE":
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID tidak valid", http.StatusBadRequest)
			return
		}

		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Task dihapus"})
				return
			}
		}
		http.Error(w, "Task tidak ditemukan", http.StatusNotFound)

	default:
		http.Error(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}

```

membuat file `middlewars.go` di folder middlewars dg isi file seperti ini:

```
package middlewares

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware mencatat setiap request yang masuk
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Selesai dalam %v", time.Since(start))
	})
}

// CORSMiddleware mengizinkan akses dari frontend React.js
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // bisa diganti domain React
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// untuk preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

```

Kemudian jalankan: 

```
go run main.go

```

## Koneksikan dengan database PostgreSQL

link download: `https://www.postgresql.org/download/`

masuk ke postgreSQl lewat terminal dengan perintah ini:

```
psql -U postgres

```

lalu buat database baru:

```
CREATE DATABASE task_manager;
\c task_manager;

```

kemudian buat tabel tasks:

```
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'ready',
    created_at TIMESTAMP DEFAULT NOW()
);

```

### Tambah library PostgreSQL di Go

kita menggunakan library pgx (modern & cepat): //jika ingin menggunakan migration lebih baik skip ini dan gunakan GORM

```
go get github.com/jackc/pgx/v5

```

gambaran structur folder nya:

```
task-manager/
│
├─ database/
│   └─ db.go
├─ handlers/
│   └─ task_handler.go
├─ models/
│   └─ models.go
├─ middlewares/
│   └─ middlewares.go
└─ main.go

```

isi file `database/db.go`:

```
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// ConnectDatabase menghubungkan ke PostgreSQL
func ConnectDatabase() {
	dsn := "postgres://postgres:12345@localhost:5432/task_manager" // ganti password sesuai PostgreSQL kamu

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}

	fmt.Println("Koneksi ke database berhasil!")
}

```

# Update main.go

```
package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager/database"
	"task-manager/handlers"
	"task-manager/middlewares"
)

func main() {
	// koneksi database
	database.ConnectDatabase()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/tasks", handlers.TaskHandler)

	handlerWithMiddleware := middlewares.LoggingMiddleware(middlewares.CORSMiddleware(mux))

	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}

```

jalankan perintah ini: `go mod tidy`

kemudian jalankan lagi: `go run main.go`

## Menggunkan GORM di project untuk migration tabels

install GORM dan driver

```
go get gorm.io/gorm
go get gorm.io/driver/postgres

```

isi file `database/db.go` jadi seperti ini:

```
package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=task_manager port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	DB = db
	fmt.Println("Koneksi ke database berhasil (GORM)")
}

```

buat file baru: `database/migrate.go` dan isi seperti ini:

```
package database

import (
	"fmt"
	"task-manager/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		panic("Gagal migrasi database: " + err.Error())
	}
	fmt.Println("Migrasi database selesai (GORM)")
}

```

update `models/models.go` seperti ini:

```
package models

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100)"`
	Role string `gorm:"type:varchar(50)"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"type:varchar(100)"`
	Description string
	Status      string `gorm:"type:varchar(20)"`
	AssignedTo  uint
}

```

tambahkan pemanggilan database.Migrate() setelah ConnectDatabase()
jadi file `main.go` kamu seperti ini:

```
package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager/database"
	"task-manager/handlers"
	"task-manager/middlewares"
)

func main() {
	// koneksi database
	database.ConnectDatabase()

	// jalankan migrasi
	database.Migrate()

	mux := http.NewServeMux()

	// route
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/tasks", handlers.TaskHandler)

	// pasang middleware (urutan penting)
	handlerWithMiddleware := middlewares.LoggingMiddleware(middlewares.CORSMiddleware(mux))

	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}

```

jalankan lagi aplikasi nya: `go run .`
