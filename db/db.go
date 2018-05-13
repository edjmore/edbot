package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("No DATABASE_URL found.")
	}
	return sql.Open("postgres", dbURL)
}
