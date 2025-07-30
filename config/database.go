package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitDB() (*sql.DB, error) {
	// Ambil connection string dari environment variable
	dbURL := "postgresql://postgres:AdibmuhA212z@db.nxhexaubbamhtfzuytiy.supabase.co:5432/postgres"
	
	var dsn string
	
	if dbURL != "" {
		// Gunakan DATABASE_URL dari environment (untuk production/Railway)
		// DATABASE_URL sudah dalam format PostgreSQL URL, jadi langsung pakai
		dsn = dbURL
		log.Println("Using DATABASE_URL from environment variable")
	} else {
		// Fallback untuk development - coba ambil individual env vars dulu  
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		
		if host != "" && port != "" && user != "" && password != "" && dbname != "" {
			// Format key-value untuk lib/pq
			dsn = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
				user, password, host, port, dbname)
			log.Println("Using individual DB environment variables")
		} else {
			// Hardcode fallback untuk development
			dsn = "user=postgres.nxhexaubbamhtfzuytiy password=AdibmuhA212z host=aws-0-ap-southeast-1.pooler.supabase.com port=5432 dbname=postgres sslmode=require"
			log.Println("Using hardcoded database connection (development mode)")
		}
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connected successfully!")
	return db, nil
}