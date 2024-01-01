package entity

const KeySearch = "search"

type BookType uint8

const (
	AudioBook BookType = 1
	TextBook  BookType = 2
)

type Books struct {
	Books      []Book
	Pagination Pagination
}

type Book struct {
	ID    uint64
	Title string
	Cover string
	Type  BookType
}
