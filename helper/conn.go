package helper

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func KoneksiDB() *sql.DB {
	err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}

		if err := db.Ping(); err != nil {
			panic(err)
		}

		fmt.Println("Koneksi Berhasil")
		return db
}