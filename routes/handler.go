package routes

import (
	"context"
	"expenses/db"

	"github.com/antonlindstrom/pgstore"
)

type Handler struct {
	ctx     context.Context
	Queries *db.Queries
	store   *pgstore.PGStore
}

func NewHandler(ctx context.Context, queries *db.Queries, store *pgstore.PGStore) *Handler {
	return &Handler{ctx: ctx, Queries: queries, store: store}
}
