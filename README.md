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
)

func main() {
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/tasks", handlers.TaskHandler)

    fmt.Println("Server berjalan di http://localhost:1001")
    log.Fatal(http.ListenAndServe(":1001", nil))
}


```

Kemudian buat folder dengan nama `handlers` dan buat file di dalamnya dengan nama `task_handlers.go` dan isi filenya seperti ini:

```
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

```

Kemudian jalankan: 

```
go run main.go

```