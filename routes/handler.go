package routes

import (
	"context"
	"expenses/db"

	"github.com/gorilla/sessions"
)

type Handler struct {
	ctx     context.Context
	Queries *db.Queries
	store   *sessions.CookieStore
}

func NewHandler(ctx context.Context, queries *db.Queries, store *sessions.CookieStore) *Handler {
	return &Handler{ctx: ctx, Queries: queries, store: store}
}
