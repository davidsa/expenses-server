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

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	queries, database := db.SetupDb()

	defer database.Close()

	ctx := context.Background()

	store, err := db.SetupSessionStore()

	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	router := mux.NewRouter()

	handler := CORSMiddleware(router)

	h := routes.NewHandler(ctx, queries, store)

	rolesRouter := router.PathPrefix("/role").Subrouter()
	rolesRouter.HandleFunc("/", h.RoleListRoute)

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", h.UserCreateRoute)
	userRouter.HandleFunc("/login", h.UserLoginRoute)
	userRouter.HandleFunc("/me", h.UserMeRoute)

	router.HandleFunc("/", HomeHandler)

	server := &http.Server{
		Handler:      handler,
		Addr:         "localhost:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer server.Close()

	log.Fatal(server.ListenAndServe())
}
