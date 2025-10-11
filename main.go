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
	database.Fresh()

	// === Router utama ===
	router := mux.NewRouter()

	// === Simpan router ke package handlers ===
	handlers.Router = router

	// === Routes dengan nama ===
	router.HandleFunc("/dashboard", handlers.DashboardHandler).
		Name("dashboard").
		Methods("GET")

	// === Static files (CSS, JS, gambar, dll) ===
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// === Middleware (logging + CORS) ===
	handlerWithMiddleware := middlewares.LoggingMiddleware(
		middlewares.CORSMiddleware(router),
	)

	// === Jalankan server ===
	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}
