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

	fmt.Println("Server berjalan di http://localhost:1001 🚀")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}
