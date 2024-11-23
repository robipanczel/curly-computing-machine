package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string

	ListBooks(ctx context.Context) ([]Book, error)
	AddBook(ctx context.Context, book BookRequest) (*primitive.ObjectID, error)
	GetBook(ctx context.Context, bookID primitive.ObjectID) (*Book, error)
	BorrowBook(ctx context.Context, bookId string, borrowerId string) (*Book, error)

	CreateAuthor(ctx context.Context, author AuthorRequest) (*primitive.ObjectID, error)
	GetAuthor(ctx context.Context, authorID primitive.ObjectID) (*Author, error)

	CreateBorrower() (*Borrower, error)
	GetBorrower() (*Borrower, error)
	BorrowedBooks() ([]Book, error)
	// TODO: Could be useful ReturnBorrowed()
}

type service struct {
	db            *mongo.Client
	booksColl     *mongo.Collection
	authorsColl   *mongo.Collection
	borrowersColl *mongo.Collection
}

var (
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	database = os.Getenv("BLUEPRINT_DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	booksColl := client.Database(database).Collection("books")
	authorsColl := client.Database(database).Collection("authors")
	borrowersColl := client.Database(database).Collection("borrowers")

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db:            client,
		booksColl:     booksColl,
		authorsColl:   authorsColl,
		borrowersColl: borrowersColl,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
