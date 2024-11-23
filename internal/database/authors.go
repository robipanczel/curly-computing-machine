package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Author struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Birthday string             `json:"birthday" bson:"birthday"`
	Email    string             `json:"email" bson:"email"`
}

type AuthorRequest struct {
	Name     string `json:"name" bson:"name"`
	Birthday string `json:"birthday" bson:"birthday"`
	Email    string `json:"email" bson:"email"`
}

func (s *service) CreateAuthor(ctx context.Context, author AuthorRequest) (*primitive.ObjectID, error) {
	authorFilter := bson.D{
		bson.E{Key: "name", Value: author.Name},
		bson.E{Key: "birthday", Value: author.Birthday},
	}
	resultByName, err := s.getAuthorByFilter(ctx, authorFilter)
	if resultByName != nil || err != nil {
		return nil, fmt.Errorf("author already exists")
	}

	emailFilter := bson.D{
		bson.E{Key: "email", Value: author.Email},
	}
	resultByEmail, err := s.getAuthorByFilter(ctx, emailFilter)
	if resultByEmail != nil || err != nil {
		return nil, fmt.Errorf("email already exists")
	}

	result, err := s.authorsColl.InsertOne(ctx, author)
	if err != nil {
		return nil, fmt.Errorf("create author: %v", err)
	}

	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (s *service) GetAuthor(ctx context.Context, authorID primitive.ObjectID) (*Author, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: authorID},
	}

	var author Author

	err := s.authorsColl.FindOne(ctx, filter).Decode(&author)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("get author: %v", err)
	}

	return &author, nil
}

func (s *service) getAuthorByFilter(ctx context.Context, filter bson.D) (*Author, error) {
	var author Author
	err := s.authorsColl.FindOne(ctx, filter).Decode(&author)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &author, nil
}
