package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/antonlindstrom/pgstore"
)

var connStr = "postgres://postgres@localhost:5432/expenses_dev?sslmode=disable"

func SetupDb() (*Queries, *sql.DB) {
	database, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return New(database), database
}

func SetupSessionStore() (*pgstore.PGStore, error) {
	key := os.Getenv("SESSION_KEY")
	return pgstore.NewPGStore(connStr, []byte(key))
}
