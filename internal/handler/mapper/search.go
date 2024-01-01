package mapper

import (
	"gharabakloo/search/internal/entity"
	"gharabakloo/search/internal/entity/dto"
)

func ToBooksResp(books entity.Books) *dto.Books {
	data := make([]dto.Book, 0, len(books.Books))
	for _, book := range books.Books {
		data = append(data, *ToBookResp(book))
	}
	return &dto.Books{
		Books:      data,
		Pagination: toPagination(books.Pagination),
	}
}

func ToBookResp(book entity.Book) *dto.Book {
	return &dto.Book{
		Title: book.Title,
		Type:  dto.BookType[book.Type],
		Cover: book.Cover,
	}
}

func toPagination(pagination entity.Pagination) dto.Pagination {
	return dto.Pagination{
		CurrentPage: pagination.CurrentPage,
		PerPage:     pagination.PerPage,
	}
}
