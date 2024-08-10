package routes

import (
	"context"
	"database/sql"
	"expenses/db"

	"github.com/antonlindstrom/pgstore"
)

type Handler struct {
	ctx     context.Context
	Queries *db.Queries
	store   *pgstore.PGStore
	db      *sql.DB
}

func NewHandler(ctx context.Context, queries *db.Queries, store *pgstore.PGStore, db *sql.DB) *Handler {
	return &Handler{ctx: ctx, Queries: queries, store: store, db: db}
}
