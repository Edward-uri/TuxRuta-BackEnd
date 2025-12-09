package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitPostgres() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""),
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", ""),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Printf("Connected to PostgreSQL database: %s", getEnv("DB_NAME", "tuxRuta"))

	db.SetMaxOpenConns(25)                  // Máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                   // Conexiones idle en el pool
	db.SetConnMaxLifetime(5 * time.Minute)  // Tiempo de vida máximo de una conexión
	db.SetConnMaxIdleTime(10 * time.Minute) // Tiempo máximo que una conexión puede estar idle
}

func GetDB() *sql.DB {
	if db == nil {
		InitPostgres()
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
