package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateBorrower(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	borrowerRequest := BorrowerRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@hotmail.com",
	}

	id, err := srv.CreateBorrower(context.Background(), borrowerRequest)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	author, err := srv.GetBorrower(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, author)

	assert.Equal(t, borrowerRequest.Name, author.Name)
	assert.Equal(t, borrowerRequest.Birthday, author.Birthday)
	assert.Equal(t, borrowerRequest.Email, author.Email)

	testcases := []struct {
		name     string
		borrower BorrowerRequest
		errMsg   string
	}{
		{
			name: "borrower already exists",
			borrower: BorrowerRequest{
				Name:     "Bober",
				Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
				Email:    "skunks@hotmail.com",
			},
			errMsg: "borrower already exists",
		},
		{
			name: "email already exists",
			borrower: BorrowerRequest{
				Name:     "Bober",
				Birthday: time.Date(1996, time.May, 18, 0, 0, 0, 0, time.UTC),
				Email:    "bober@hotmail.com",
			},
			errMsg: "email already exists",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			id, err := srv.CreateBorrower(context.Background(), testcase.borrower)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), testcase.errMsg)
			assert.Nil(t, id)
		})
	}
}

func TestGetBorrower(t *testing.T) {
	srv := New()

	err := srv.(*service).deleteColls(context.Background())
	assert.NoError(t, err)

	borrowerRequest := BorrowerRequest{
		Name:     "Bober",
		Birthday: time.Date(1996, time.May, 17, 0, 0, 0, 0, time.UTC),
		Email:    "bober@hotmail.com",
	}

	id, err := srv.CreateBorrower(context.Background(), borrowerRequest)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	author, err := srv.GetBorrower(context.Background(), *id)
	assert.NoError(t, err)
	assert.NotNil(t, author)

	assert.Equal(t, borrowerRequest.Name, author.Name)
	assert.Equal(t, borrowerRequest.Birthday, author.Birthday)
	assert.Equal(t, borrowerRequest.Email, author.Email)

	noBorrower, err := srv.GetBorrower(context.Background(), primitive.NewObjectID())
	assert.Nil(t, noBorrower)
	assert.NoError(t, err)

	noIdBorrower, err := srv.GetBorrower(context.Background(), primitive.ObjectID{})
	assert.Nil(t, noIdBorrower)
	assert.NoError(t, err)
}
