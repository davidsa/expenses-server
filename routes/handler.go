package routes

import (
	"context"
	"expenses/db"
)

type Handler struct {
	ctx     context.Context
	Queries *db.Queries
}

func NewHandler(ctx context.Context, queries *db.Queries) *Handler {
	return &Handler{ctx: ctx, Queries: queries}
}
