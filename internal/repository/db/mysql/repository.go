package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/pkg/myerr"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Search(ctx context.Context, key string, p entity.Pagination) (*entity.Books, error) {
	likeKey := fmt.Sprintf("%%%s%%", key)
	query := `
		SELECT id, title, cover, type 
		FROM books 
		WHERE LOWER(title) LIKE ? 
		ORDER BY title 
		LIMIT ? OFFSET ?;
	`

	rows, err := r.db.QueryContext(ctx, query, likeKey, p.PerPage, p.GetOffset())
	if err != nil {
		return nil, myerr.Errorf(err)
	}
	defer func() { _ = rows.Close() }()

	books := make([]entity.Book, 0)
	for rows.Next() {
		var book entity.Book
		if err = rows.Scan(&book.ID, &book.Title, &book.Cover, &book.Type); err != nil {
			return nil, myerr.Errorf(err)
		}

		books = append(books, book)
	}
	return &entity.Books{
		Books:      books,
		Pagination: p,
	}, nil
}
