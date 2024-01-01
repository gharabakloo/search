package dto

import "gharabakloo/search/internal/entity"

var BookType = map[entity.BookType]string{
	entity.AudioBook: "audiobook",
	entity.TextBook:  "textbook",
}

type Books struct {
	Books      []Book     `json:"books"`
	Pagination Pagination `json:"pagination"`
}

type Book struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Cover string `json:"cover"`
}

type Pagination struct {
	CurrentPage uint64 `json:"current_page"`
	PerPage     uint64 `json:"per_page"`
}
