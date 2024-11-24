package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Borrower struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id"`
	Name     string               `json:"name" bson:"name"`
	Birthday time.Time            `json:"birthday" bson:"birthday"`
	Email    string               `json:"email" bson:"email"`
	Books    []primitive.ObjectID `json:"books" bson:"books"`
}

func (b *Borrower) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type BorrowerRequest struct {
	Name     string    `json:"name" bson:"name"`
	Birthday time.Time `json:"birthday" bson:"birthday"`
	Email    string    `json:"email" bson:"email"`
}

func (b *BorrowerRequest) Bind(r *http.Request) error {
	if b.Birthday.IsZero() {
		return fmt.Errorf("birthday is required")
	}

	if b.Email == "" {
		return fmt.Errorf("email is required")
	}

	if b.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
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

func (s *service) BorrowedBooks(ctx context.Context, borrowerID primitive.ObjectID) ([]Book, error) {
	borrower, err := s.GetBorrower(ctx, borrowerID)
	if err != nil {
		return nil, fmt.Errorf("get borrower: %v", err)
	}
	if borrower == nil {
		return nil, fmt.Errorf("borrower doesn't exist")
	}

	if len(borrower.Books) == 0 {
		return []Book{}, nil
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": borrower.Books,
		},
	}

	cursor, err := s.booksColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find books: %v", err)
	}
	defer cursor.Close(ctx)

	var books []Book
	err = cursor.All(ctx, &books)
	if err != nil {
		return nil, fmt.Errorf("decode books: %v", err)
	}

	return books, nil
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

func (s *service) borrowBookByUser(ctx context.Context, borrowerID primitive.ObjectID, bookID primitive.ObjectID) error {
	filter := bson.M{
		"_id": borrowerID,
	}

	update := bson.M{
		"$push": bson.M{
			"books": bookID,
		},
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedBorrower Borrower
	err := s.borrowersColl.FindOneAndUpdate(ctx, filter, update, opt).Decode(&updatedBorrower)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("borrower doesn't exist")
		}
		return fmt.Errorf("find and update: %v", err)
	}

	return nil
}
