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

	// === Routes dashboard ===
	router.HandleFunc("/dashboard", handlers.DashboardHandler).
		Name("dashboard").
		Methods("GET")

	// users
	router.HandleFunc("/users", handlers.UsersHandler).
		Name("users").
		Methods("GET")

	router.HandleFunc("/users/create", handlers.CreateUserHandler).
		Name("users.create").
		Methods("GET")

	router.HandleFunc("/users/store", handlers.StoreUserHandler).
		Name("users.store").
		Methods("POST")

	router.HandleFunc("/users/edit/{id}", handlers.EditUserHandler).
		Name("users.edit").
		Methods("GET")

	router.HandleFunc("/users/update/{id}", handlers.UpdateUserHandler).
		Name("users.update").
		Methods("POST")
	
	router.HandleFunc("/users/delete/{id}", handlers.DeleteUserHandler).
		Name("users.delete").
		Methods("POST")

	// === Simpan router ke package handlers (setelah semua routes terdaftar) ===
	handlers.Router = router

	// === Static files ===
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// === Middleware ===
	handlerWithMiddleware := middlewares.LoggingMiddleware(
		middlewares.CORSMiddleware(router),
	)

	// === Jalankan server ===
	fmt.Println("Server berjalan di http://localhost:1001")
	log.Fatal(http.ListenAndServe(":1001", handlerWithMiddleware))
}
