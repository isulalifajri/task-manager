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
	dsn := "postgres://postgres:postgres@localhost:5432/task_manager" // ganti password sesuai PostgreSQL kamu

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
