package redis

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/pkg/myerr"
)

type Repository struct {
	client entity.RedisClient
}

func New(client entity.RedisClient) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Add(ctx context.Context, key string, books *entity.Books) error {
	member, err := json.Marshal(books)
	if err != nil {
		return myerr.Errorf(err)
	}

	z := redis.Z{
		Score:  float64(books.Pagination.CurrentPage),
		Member: string(member),
	}
	return myerr.Errorf(r.client.ZAdd(ctx, key, z).Err())
}

func (r *Repository) Get(ctx context.Context, key string, cacheRange entity.Range) ([]*entity.Books, error) {
	result, err := r.client.ZRange(ctx, key, cacheRange.Start, cacheRange.End).Result()
	if err != nil {
		return nil, myerr.Errorf(err)
	}

	booksPages := make([]*entity.Books, len(result))
	for _, data := range result {
		books := new(entity.Books)
		if err = json.Unmarshal([]byte(data), books); err != nil {
			return nil, myerr.Errorf(err)
		}

		booksPages = append(booksPages, books)
	}
	return booksPages, nil
}
