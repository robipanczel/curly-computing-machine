package database

type Author struct {
	ID         string
	AuthorName string
	BirthDay   string
	Email      string
}

func (s *service) CreateAuthor() (*Author, error) {
	return nil, nil
}

func (s *service) GetAuthor() (*Author, error) {
	return nil, nil
}
