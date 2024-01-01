package db

import (
	"context"

	"gharabakloo/search/internal/entity"
)

type Repository interface {
	Search(ctx context.Context, key string, p entity.Pagination) (*entity.Books, error)
}
