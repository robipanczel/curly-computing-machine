package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	AuthorID    primitive.ObjectID `json:"author_id" bson:"author_id"`
	Genres      []string           `json:"genres" bson:"genres"`
	Available   bool               `json:"available" bson:"available"`
}

type BookRequest struct {
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	AuthorID    primitive.ObjectID `json:"author_id" bson:"author_id"`
	Genres      []string           `json:"genres" bson:"genres"`
	Available   bool               `json:"available" bson:"available"`
}

func (s *service) ListBooks(ctx context.Context) ([]Book, error) {
	filter := bson.D{}

	var books []Book

	curs, err := s.booksColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find books: %v", err)
	}

	err = curs.All(ctx, &books)
	if err != nil {
		return nil, fmt.Errorf("decode all books: %v", err)
	}

	return books, nil
}

func (s *service) AddBook(ctx context.Context, book BookRequest) (*primitive.ObjectID, error) {
	author, err := s.GetAuthor(ctx, book.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get author: %v", err)
	}
	if author == nil {
		return nil, fmt.Errorf("author doesn't exists")
	}

	bookFilter := bson.D{
		bson.E{Key: "title", Value: book.Title},
		bson.E{Key: "author_id", Value: book.AuthorID},
	}
	bookExists, err := s.getBookByFilter(ctx, bookFilter)
	if err != nil {
		return nil, fmt.Errorf("book validating: %v", err)
	}
	if bookExists != nil {
		return nil, fmt.Errorf("book already exists")
	}

	result, err := s.booksColl.InsertOne(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("failed to create new book: %v", err)
	}

	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (s *service) GetBook(ctx context.Context, bookID primitive.ObjectID) (*Book, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: bookID},
	}

	var book Book

	err := s.booksColl.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("get book: %v", err)
	}

	return &book, nil
}

func (s *service) BorrowBook(ctx context.Context, bookID primitive.ObjectID, borrowerID primitive.ObjectID) error {
	book, err := s.GetBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("get book: %v", err)
	}

	if book == nil {
		return fmt.Errorf("book doesn't exist")
	}

	borrower, err := s.GetBorrower(ctx, borrowerID)
	if err != nil {
		return fmt.Errorf("get borrower %v", err)
	}

	if borrower == nil {
		return fmt.Errorf("borrower doesn't exist")
	}

	if !book.Available {
		return fmt.Errorf("book isn't available")
	}

	update := bson.M{
		"$set": bson.M{
			"available": false,
		},
	}

	_, err = s.booksColl.UpdateByID(ctx, book.ID, update)
	if err != nil {
		return fmt.Errorf("book available update: %v", err)
	}

	err = s.borrowBookByUser(ctx, borrowerID, bookID)
	if err != nil {
		return fmt.Errorf("borrow book by user: %v", err)
	}

	return nil
}

func (s *service) getBookByFilter(ctx context.Context, filter bson.D) (*Book, error) {
	var Book Book
	err := s.booksColl.FindOne(ctx, filter).Decode(&Book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &Book, nil
}
