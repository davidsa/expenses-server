package routes

import (
	"database/sql"
	"encoding/json"
	"expenses/db"
	"expenses/utils"
	"fmt"
	"log"
	"net/http"
)

type CreateUserPayload struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	RoleID   int32  `json:"role_id"`
}

func (h Handler) UserCreateRoute(w http.ResponseWriter, r *http.Request) {

	var p CreateUserPayload

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password_hash, err := utils.HashPassword(p.Password)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not hash the password", http.StatusBadRequest)
		return
	}

	user := db.CreateUserParams{
		Email:        p.Email,
		Name:         p.Name,
		Lastname:     p.Lastname,
		PasswordHash: password_hash,
		RoleID:       sql.NullInt32{Int32: 2, Valid: true},
	}

	result, err := h.Queries.CreateUser(h.ctx, user)

	response := CreateUserResponse{
		ID:       result.ID,
		Name:     result.Name,
		Lastname: result.Lastname,
		Email:    result.Email,
		RoleID:   result.RoleID.Int32,
	}

	if err != nil {
		fmt.Println(err.Error())

		if encodeErr := json.NewEncoder(w).Encode(utils.JsonError{Error: err.Error()}); encodeErr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
