package database

type Borrower struct {
	ID       string
	Name     string
	BirthDay string
	Email    string
}

func (s *service) CreateBorrower() (*Borrower, error) {
	return nil, nil
}

func (s *service) GetBorrower() (*Borrower, error) {
	return nil, nil
}

func (s *service) BorrowedBooks() ([]Book, error) {
	return nil, nil
}
