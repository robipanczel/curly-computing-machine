package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Borrower struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Birthday string             `json:"birthday" bson:"birthday"`
	Email    string             `json:"email" bson:"email"`
}

type BorrowerRequest struct {
	Name     string `json:"name" bson:"name"`
	Birthday string `json:"birthday" bson:"birthday"`
	Email    string `json:"email" bson:"email"`
}

func (s *service) CreateBorrower(ctx context.Context, borrower BorrowerRequest) (*primitive.ObjectID, error) {
	borrowerFilter := bson.D{
		bson.E{Key: "name", Value: borrower.Name},
		bson.E{Key: "birthday", Value: borrower.Birthday},
	}
	resultByName, err := s.getBorrowerByFilter(ctx, borrowerFilter)
	if resultByName != nil || err != nil {
		return nil, fmt.Errorf("borrower already exists")
	}

	emailFilter := bson.D{
		bson.E{Key: "email", Value: borrower.Email},
	}
	resultByEmail, err := s.getBorrowerByFilter(ctx, emailFilter)
	if resultByEmail != nil || err != nil {
		return nil, fmt.Errorf("email already exists")
	}

	result, err := s.borrowersColl.InsertOne(ctx, borrower)
	if err != nil {
		return nil, fmt.Errorf("create borrower: %v", err)
	}

	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (s *service) GetBorrower(ctx context.Context, borrowerID primitive.ObjectID) (*Borrower, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: borrowerID},
	}

	var borrower Borrower

	err := s.borrowersColl.FindOne(ctx, filter).Decode(&borrower)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("get borrower: %v", err)
	}

	return &borrower, nil
}

func (s *service) BorrowedBooks(ctx context.Context) ([]Borrower, error) {
	return nil, nil
}

func (s *service) getBorrowerByFilter(ctx context.Context, filter bson.D) (*Borrower, error) {
	var borrower Borrower
	err := s.borrowersColl.FindOne(ctx, filter).Decode(&borrower)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &borrower, nil
}
