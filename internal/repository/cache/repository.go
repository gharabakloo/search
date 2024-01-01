package cache

import (
	"context"

	"gharabakloo/search/internal/entity"
)

type Repository interface {
	Add(ctx context.Context, key string, books *entity.Books) error
	Get(ctx context.Context, key string, cacheRange entity.Range) ([]*entity.Books, error)
}
