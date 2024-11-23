package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Borrower struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Birthday string             `json:"birthday" bson:"birthday"`
	Email    string             `json:"email" bson:"email"`
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
