package mysql

import (
	"context"
	"fmt"
	"gharabakloo/search/internal/entity"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestSearchWithoutError
func TestSearchWithoutError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}

	key := "book"
	p := entity.Pagination{
		CurrentPage: 1,
		PerPage:     100,
	}

	mockBooks := []entity.Book{
		{
			ID:    1,
			Title: "book1",
			Cover: "cover1",
			Type:  2,
		},
		{
			ID:    2,
			Title: "book2",
			Cover: "cover2",
			Type:  2,
		},
		{
			ID:    3,
			Title: "book12",
			Cover: "cover12",
			Type:  1,
		},
	}

	exceptedBooks := &entity.Books{
		Books:      mockBooks,
		Pagination: p,
	}

	rows := sqlmock.NewRows([]string{"id", "title", "cover", "type"}).
		AddRow(mockBooks[0].ID, mockBooks[0].Title, mockBooks[0].Cover, mockBooks[0].Type).
		AddRow(mockBooks[1].ID, mockBooks[1].Title, mockBooks[1].Cover, mockBooks[1].Type).
		AddRow(mockBooks[2].ID, mockBooks[2].Title, mockBooks[2].Cover, mockBooks[2].Type)

	query := `
		SELECT id, title, cover, type 
		FROM books 
		WHERE LOWER(title) LIKE ? 
		ORDER BY title 
		LIMIT ? OFFSET ?;
	`

	likeKey := fmt.Sprintf("%%%s%%", key)
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(likeKey, p.PerPage, p.GetOffset()).
		WillReturnRows(rows)

	repo := Repository{db: db}
	actualBooks, err := repo.Search(context.Background(), key, p)

	assert.Nil(t, err)
	assert.EqualValues(t, exceptedBooks, actualBooks)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err.Error())
	}
}
