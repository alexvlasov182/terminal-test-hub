package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgres creates new Postgres database connection
//
// Parameters:
//   - dsn: Data Source Name for the PostgreSQL connection"
//     "postgres://username:password@localhost:5432/dbname?sslmode=disable"
//
// Returns:
//   - *sqlx.DB: Pointer to the sqlx.DB instance representing the database connection
//   - error: Error if the connection could not be established, otherwise nil
func NewPostgres(dsn string) (*sqlx.DB, error) {
	// connect to postgres database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot to connect to the PostgreSQL: %w", err)
	}

	// check if the database is reachable
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping the PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(25)                 // maximum number of open connections to the database
	db.SetMaxIdleConns(5)                  // 5 connections can be idle
	db.SetConnMaxLifetime(5 * time.Minute) // connection can be reused for 5 minutes

	log.Println("✅ Connected to PostgreSQL database successfully")
	return db, nil

}
