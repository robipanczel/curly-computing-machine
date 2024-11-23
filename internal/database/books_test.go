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

	bookRequest := BookRequest{
		Title:       "Hobbit",
		Description: "The Hobbit is set in Middle-earth",
		AuthorID:    *authorID,
		Genres:      []string{"fantasy"},
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

	emptyBooks, err := srv.ListBooks(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(emptyBooks))

	bookRequest := BookRequest{
		Title:       "Hobbit",
		Description: "The Hobbit is set in Middle-earth",
		AuthorID:    *authorID,
		Genres:      []string{"fantasy"},
	}
	bookID, err := srv.AddBook(context.Background(), bookRequest)
	assert.NoError(t, err)
	assert.NotNil(t, bookID)

	books, err := srv.ListBooks(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(books))
	assert.Equal(t, bookRequest.Title, books[0].Title)
}
