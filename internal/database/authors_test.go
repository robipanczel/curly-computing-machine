package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateAuthor(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	authorRequest := AuthorRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@author.com",
	}

	id, err := srv.CreateAuthor(context.Background(), authorRequest)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	author, err := srv.GetAuthor(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, author)

	assert.Equal(t, authorRequest.Name, author.Name)
	assert.Equal(t, authorRequest.Birthday, author.Birthday)
	assert.Equal(t, authorRequest.Email, author.Email)

	testcases := []struct {
		name   string
		author AuthorRequest
		errMsg string
	}{
		{
			name: "author already exists",
			author: AuthorRequest{
				Name:     "Bober",
				Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
				Email:    "pingvin@author.com",
			},
			errMsg: "author already exists",
		},
		{
			name: "email already exists",
			author: AuthorRequest{
				Name:     "Skunks",
				Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
				Email:    "bober@author.com",
			},
			errMsg: "email already exists",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			id, err := srv.CreateAuthor(context.Background(), testcase.author)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), testcase.errMsg)
			assert.Nil(t, id)
		})
	}
}

func TestGetAuthor(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	authorRequest := AuthorRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@author.com",
	}

	id, err := srv.CreateAuthor(context.Background(), authorRequest)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	author, err := srv.GetAuthor(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, author)

	assert.Equal(t, authorRequest.Name, author.Name)
	assert.Equal(t, authorRequest.Birthday, author.Birthday)
	assert.Equal(t, authorRequest.Email, author.Email)

	noAuthor, err := srv.GetAuthor(context.Background(), primitive.NewObjectID())
	assert.Nil(t, noAuthor)
	assert.NoError(t, err)

	noIdAuthor, err := srv.GetAuthor(context.Background(), primitive.ObjectID{})
	assert.Nil(t, noIdAuthor)
	assert.NoError(t, err)
}
