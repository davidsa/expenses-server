package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h Handler) RoleListHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := h.Queries.ListRoles(h.ctx)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(roles)
}
