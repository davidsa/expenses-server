package main

import (
	"context"
	"encoding/gob"
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
	fmt.Fprintf(w, "ok")
}

func init() {
	gob.Register(routes.UserCreateResponse{})
}

func main() {

	queries, database := db.SetupDb()

	defer database.Close()

	router := mux.NewRouter()

	ctx := context.Background()

	store, err := db.SetupSessionStore()

	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	h := routes.NewHandler(ctx, queries, store)

	rolesRouter := router.PathPrefix("/role").Subrouter()
	rolesRouter.HandleFunc("/", h.RoleListRoute)

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", h.UserCreateRoute)
	userRouter.HandleFunc("/login", h.UserLoginRoute)
	userRouter.HandleFunc("/me", h.UserMeRoute)

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
