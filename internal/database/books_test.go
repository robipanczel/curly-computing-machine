package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddBook(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	authorRequest := AuthorRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@author.com",
	}
	authorID, err := srv.CreateAuthor(context.Background(), authorRequest)
	assert.NoError(t, err)
	assert.NotNil(t, authorID)

	t.Run("should create book", func(t *testing.T) {
		bookRequest := BookRequest{
			Title:       "Hobbit",
			Description: "The Hobbit is set in Middle-earth",
			AuthorID:    *authorID,
			Genres:      []string{"fantasy"},
			Available:   true,
		}
		bookID, err := srv.AddBook(context.Background(), bookRequest)
		assert.NoError(t, err)
		assert.NotNil(t, bookID)

		book, err := srv.GetBook(context.Background(), *bookID)
		assert.NoError(t, err)
		assert.NotNil(t, book)

		assert.Equal(t, bookRequest.AuthorID.Hex(), book.AuthorID.Hex())
		assert.Equal(t, bookRequest.Description, book.Description)
		assert.Equal(t, bookRequest.Genres, book.Genres)
		assert.Equal(t, bookRequest.Title, book.Title)
	})

	testcases := []struct {
		name   string
		book   BookRequest
		errMsg string
	}{
		{
			name: "author doesn't exists",
			book: BookRequest{
				Title:       "Hoho",
				Description: "This is a story about winter",
				AuthorID:    primitive.NewObjectID(),
				Genres:      []string{},
				Available:   true,
			},
			errMsg: "author doesn't exists",
		},
		{
			name: "book already exists",
			book: BookRequest{
				Title:       "Hobbit",
				Description: "The Hobbit is set in Middle-earth",
				AuthorID:    *authorID,
				Genres:      []string{"fantasy"},
				Available:   true,
			},
			errMsg: "book already exists",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			id, err := srv.AddBook(context.Background(), testcase.book)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), testcase.errMsg)
			assert.Nil(t, id)
		})
	}
}

func TestListBooks(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	authorRequest := AuthorRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@author.com",
	}
	authorID, err := srv.CreateAuthor(context.Background(), authorRequest)
	assert.NoError(t, err)
	assert.NotNil(t, authorID)

	t.Run("should list no books if db empty", func(t *testing.T) {
		emptyBooks, err := srv.ListBooks(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(emptyBooks))
	})

	t.Run("should list books", func(t *testing.T) {
		bookRequest := BookRequest{
			Title:       "Hobbit",
			Description: "The Hobbit is set in Middle-earth",
			AuthorID:    *authorID,
			Genres:      []string{"fantasy"},
			Available:   true,
		}
		bookID, err := srv.AddBook(context.Background(), bookRequest)
		assert.NoError(t, err)
		assert.NotNil(t, bookID)

		books, err := srv.ListBooks(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(books))
		assert.Equal(t, bookRequest.Title, books[0].Title)
	})

}

func TestBorrowBook(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	borrowerRequest := BorrowerRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@hotmail.com",
	}

	borrowerID, err := srv.CreateBorrower(context.Background(), borrowerRequest)
	assert.NoError(t, err)
	assert.NotNil(t, borrowerID)

	authorRequest := AuthorRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@author.com",
	}

	authorID, err := srv.CreateAuthor(context.Background(), authorRequest)
	assert.NoError(t, err)
	assert.NotNil(t, authorID)

	bookRequest := BookRequest{
		Title:       "Hobbit",
		Description: "The Hobbit is set in Middle-earth",
		AuthorID:    *authorID,
		Genres:      []string{"fantasy"},
		Available:   true,
	}
	bookID, err := srv.AddBook(context.Background(), bookRequest)
	assert.NoError(t, err)
	assert.NotNil(t, bookID)

	t.Run("should borrow a book, set book to not available and borrower have book", func(t *testing.T) {
		err := srv.BorrowBook(context.Background(), *bookID, *borrowerID)
		assert.NoError(t, err)

		book, err := srv.GetBook(context.Background(), *bookID)
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, false, book.Available)

		borrower, err := srv.GetBorrower(context.Background(), *borrowerID)
		assert.NoError(t, err)
		assert.NotNil(t, borrower)
		assert.Equal(t, *bookID, borrower.Books[0])
	})

	testcases := []struct {
		name       string
		bookID     primitive.ObjectID
		borrowerID primitive.ObjectID
		errMsg     string
	}{
		{
			name:       "should not borrow a book if book doesn't exist",
			bookID:     primitive.NewObjectID(),
			borrowerID: *borrowerID,
			errMsg:     "book doesn't exist",
		},
		{
			name:       "should not borrow a book if borrower doesn't exist",
			bookID:     *bookID,
			borrowerID: primitive.NewObjectID(),
			errMsg:     "borrower doesn't exist",
		},
		{
			name:       "should not borrow a book if book isn't available",
			bookID:     *bookID,
			borrowerID: *borrowerID,
			errMsg:     "book isn't available",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			err := srv.BorrowBook(context.Background(), testcase.bookID, testcase.borrowerID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), testcase.errMsg)
		})
	}

}
