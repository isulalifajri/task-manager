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
