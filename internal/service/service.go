package service

import (
	"context"

	"gharabakloo/search/internal/entity"
)

type SearchService interface {
	Search(ctx context.Context, key string, page entity.Page) (*entity.Books, error)
}
