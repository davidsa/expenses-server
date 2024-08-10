package routes

import (
	"database/sql"
	"encoding/json"
	"expenses/db"
	"fmt"
	"net/http"
)

type ListUserGroupsResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type GroupCreatePayload struct {
	Name string `json:"name"`
}

type GroupCreateResponse struct {
	Ok bool `json:"ok"`
}

func (h Handler) GroupListRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")

	if user, ok := session.Values["user"].(UserCreateResponse); ok {

		fmt.Println("User id", user.ID)
		result, err := h.Queries.ListUserGroups(h.ctx, user.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println("Result", result)

		if result == nil {
			json.NewEncoder(w).Encode([]db.ListUserGroupsRow{})
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(result)
	}
}

func (h Handler) GroupCreateRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")

	var p GroupCreatePayload

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: guard auth/unauth routes
	if user, ok := session.Values["user"].(UserCreateResponse); ok {

		tx, err := h.db.Begin()

		defer func() {

			if err != nil {

				if rbErr := tx.Rollback(); rbErr != nil {
					fmt.Printf("Failed to follback transaction: %v", rbErr)
				}
				return
			}

			if commitErr := tx.Commit(); commitErr != nil {
				fmt.Printf("Failed to commit transaction: %v", commitErr)
			}

		}()

		qtx := h.Queries.WithTx(tx)

		createResult, err := qtx.CreateGroup(h.ctx, p.Name)

		if err != nil {
			err = fmt.Errorf("Failed to create user: %w", err)
		}

		params := db.AddUserToGroupParams{
			GroupID: createResult.ID,
			UserID:  user.ID,
			IsAdmin: sql.NullBool{Bool: true, Valid: true},
		}

		if err := qtx.AddUserToGroup(h.ctx, params); err != nil {
			err = fmt.Errorf("Failed to associate user to group: %w", err)
		}

		result := GroupCreateResponse{
			Ok: true,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(result)
	}
}
