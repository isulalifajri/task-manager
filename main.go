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
	// === Setup Database ===
	database.ConnectDatabase()
	database.Migrate()
	database.Seed()

	// === Router utama ===
	router := mux.NewRouter()

	// === Route Dashboard HTML ===
	router.HandleFunc("/dashboard", handlers.DashboardHandler).Methods("GET")

	// === Static files (CSS, JS, gambar, dll) ===
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// === Pasang middleware ===
	handlerWithMiddleware := middlewares.LoggingMiddleware(
		middlewares.CORSMiddleware(router),
	)

	// === Jalankan server ===
	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}
