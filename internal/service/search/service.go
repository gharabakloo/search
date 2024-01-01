package search

import (
	"context"
	"strings"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/internal/repository/cache"
	"gharabakloo/search/internal/repository/db"
	"gharabakloo/search/pkg/myerr"
)

type Service struct {
	db    db.Repository
	cache cache.Repository
}

func New(db db.Repository, cache cache.Repository) *Service {
	return &Service{
		db:    db,
		cache: cache,
	}
}

// Search implement search key from database and cache.
// First, it searches key in cache and if it completely exists in cache,
// return response, but if not exist in cache, or it exists partially in cache,
// it queries to database for fetch complete result and save response to cache
// and send it to client.
//
// For cache result in redis, I use from sorted set in redis.
// In such a way that I keep a sorted set in Redis for each searched key.
// For handling pagination in cache, I consider a constant number for page size (for example 100)
// and save each page as a member with page number in sorted set.
// After than for any request with any pagination, first compute cache pagination,
// and then fetch data from cache.
func (s *Service) Search(ctx context.Context, key string, page entity.Page) (*entity.Books, error) {
	key = strings.ToLower(key)
	pagination := page.Parse()
	cacheRange := entity.NewRange(pagination.GetFrom(), pagination.GetTo())
	booksPages, err := s.cache.Get(ctx, key, cacheRange)
	if err != nil {
		return nil, myerr.Errorf(err)
	}

	// key doesn't exist in cache
	if len(booksPages) == 0 {
		arrBooks, err := s.searchAndCacheRange(ctx, key, cacheRange)
		if err != nil {
			return nil, myerr.Errorf(err)
		}
		return getBooks(arrBooks, pagination), nil
	}

	// complete result exist in cache
	if len(booksPages) == cacheRange.Diff() {
		arrBooks := make([]entity.Book, 0)
		for _, books := range booksPages {
			arrBooks = append(arrBooks, books.Books...)
		}
		return getBooks(arrBooks, pagination), nil
	}

	// partial result exist in cache
	arrBooks, err := s.searchAndCachePartial(ctx, booksPages, key, cacheRange)
	if err != nil {
		return nil, myerr.Errorf(err)
	}
	return getBooks(arrBooks, pagination), nil
}

func (s *Service) searchAndCacheRange(ctx context.Context, key string, r entity.Range) ([]entity.Book, error) {
	dbPagination := entity.Pagination{
		PerPage: entity.CachePageSize,
	}

	arrBooks := make([]entity.Book, 0)
	for i := r.Start; i <= r.End; i++ {
		dbPagination.CurrentPage = uint64(i)
		books, err := s.searchAndCache(ctx, key, dbPagination)
		if err != nil {
			return nil, myerr.Errorf(err)
		}

		if len(books.Books) == 0 {
			break
		}

		arrBooks = append(arrBooks, books.Books...)
	}
	return arrBooks, nil
}

func (s *Service) searchAndCache(ctx context.Context, key string, p entity.Pagination) (*entity.Books, error) {
	books, err := s.db.Search(ctx, key, p)
	if err != nil {
		return nil, myerr.Errorf(err)
	}

	if len(books.Books) == 0 {
		return books, nil
	}

	if err = s.cache.Add(ctx, key, books); err != nil {
		return nil, myerr.Errorf(err)
	}
	return books, nil
}

func (s *Service) searchAndCachePartial(ctx context.Context, booksPages []*entity.Books, key string, r entity.Range) (
	[]entity.Book, error) {
	dbPagination := entity.Pagination{
		PerPage: entity.CachePageSize,
	}

	idx := 0
	length := len(booksPages)
	arrBooks := make([]entity.Book, 0)
	for i := r.Start; i <= r.End; i++ {
		if idx < length && booksPages[idx].Pagination.CurrentPage == uint64(i) {
			arrBooks = append(arrBooks, booksPages[idx].Books...)
			idx++
			continue
		}

		dbPagination.CurrentPage = uint64(i)
		books, err := s.searchAndCache(ctx, key, dbPagination)
		if err != nil {
			return nil, myerr.Errorf(err)
		}

		if len(books.Books) == 0 {
			break
		}

		arrBooks = append(arrBooks, books.Books...)
	}
	return arrBooks, nil
}

func getBooks(arrBooks []entity.Book, p entity.Pagination) *entity.Books {
	from := (p.GetFrom() - 1) % entity.CachePageSize
	to := from + p.PerPage
	l := uint64(len(arrBooks))
	if l < to {
		to = l
	}
	return &entity.Books{
		Books:      arrBooks[from:to],
		Pagination: p,
	}
}
