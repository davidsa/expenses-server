package main

import (
	"context"
	"encoding/gob"
	"expenses/db"
	"expenses/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!!")
}

func initGob() {
	gob.Register(routes.UserCreateResponse{})
}

func main() {

	queries, db := db.SetupDb()

	defer db.Close()

	router := mux.NewRouter()

	ctx := context.Background()

	initGob()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

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
