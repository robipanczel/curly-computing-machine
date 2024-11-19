package database

type Book struct {
	ID         string
	BookName   string
	AuthorID   string
	Genres     []string
	BorrowerID string
}

func (s *service) ListBooks() ([]Book, error) {
	return nil, nil
}

func (s *service) AddBook() (*Book, error) {
	return nil, nil
}

func (s *service) BorrowBook() (*Book, error) {
	return nil, nil
}
