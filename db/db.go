package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var conn *sql.DB

func init() {
	var err error
	conn, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	createSchema()
}

func createSchema() {
	q := `
    CREATE TABLE IF NOT EXISTS messages(
      created_at INTEGER NOT NULL,
      group_id TEXT NOT NULL,
      sender_id TEXT NOT NULL,
      text TEXT NOT NULL
    )`
	if _, err := conn.Exec(q); err != nil {
		log.Fatal(err)
	}
}

func Conn() *sql.DB {
	return conn
}
