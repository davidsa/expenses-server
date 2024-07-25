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

type UserCreatePayload struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
	RoleID   int32  `json:"role_id"`
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Ok bool `json:"ok"`
}

func (h Handler) UserCreateRoute(w http.ResponseWriter, r *http.Request) {

	var p UserCreatePayload

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

	response := UserCreateResponse{
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
			return
		}
	}

	session, _ := h.store.Get(r, "session")

	session.Values["user"] = response

	sessionErr := session.Save(r, w)

	if sessionErr != nil {
		http.Error(w, sessionErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func (h Handler) UserLoginRoute(w http.ResponseWriter, r *http.Request) {

	var p UserLoginPayload

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.Queries.FindUserByEmail(h.ctx, p.Email)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(utils.JsonError{Error: "unauthorized"})
		return
	}

	match := utils.ComparePasswords(user.PasswordHash, p.Password)

	if !match {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(utils.JsonError{Error: "unauthorized"})
		return
	}

	user_session := UserCreateResponse{
		ID:       user.ID,
		Name:     user.Name,
		Lastname: user.Lastname,
		Email:    user.Email,
		RoleID:   user.RoleID.Int32,
	}

	session, _ := h.store.Get(r, "session")

	session.Values["user"] = user_session

	sessionErr := session.Save(r, w)

	if sessionErr != nil {
		http.Error(w, sessionErr.Error(), http.StatusInternalServerError)
		return
	}

	response := UserLoginResponse{
		Ok: true,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func (h Handler) UserMeRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")

	if user, ok := session.Values["user"].(UserCreateResponse); ok {

		response := UserCreateResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Lastname: user.Lastname,
			RoleID:   user.RoleID,
		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)

	} else {
		json.NewEncoder(w).Encode(utils.JsonError{Error: "No session data"})
	}
}
