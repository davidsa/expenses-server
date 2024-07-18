package main

import (
	"context"
	"expenses/db"
	"expenses/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!!")
}

func main() {

	queries, db := db.SetupDb()

	defer db.Close()

	router := mux.NewRouter()

	ctx := context.Background()
	h := routes.NewHandler(ctx, queries)

	rolesRouter := router.PathPrefix("/role").Subrouter()
	rolesRouter.HandleFunc("/", h.RoleListHandler)

	router.HandleFunc("/", HomeHandler)

	server := &http.Server{
		Handler:      router,
		Addr:         "localhost:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer server.Close()

	log.Fatal(server.ListenAndServe())
}
