package main

import (
    "fmt"
    "log"
    "net/http"
    "task-manager/database"
    "task-manager/handlers"
    "task-manager/middlewares"

    "github.com/gorilla/mux"
)

func main() {
    // koneksi database
    database.ConnectDatabase()
    database.Migrate()
    database.Seed()

    router := mux.NewRouter()

    // route
    router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
    router.HandleFunc("/tasks", handlers.TaskHandler).Methods("GET", "POST")
    router.HandleFunc("/tasks/{id}", handlers.TaskHandler).Methods("PUT", "DELETE")

    // pasang middleware
    handlerWithMiddleware := middlewares.LoggingMiddleware(middlewares.CORSMiddleware(router))

    fmt.Println("Server berjalan di http://localhost:1001")
    log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}
