package db

import (
	"database/sql"
	"log"
)

func SetupDb() (*Queries, *sql.DB) {
	connStr := "postgres://postgres@localhost:5432/expenses_dev?sslmode=disable"

	database, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return New(database), database
}
