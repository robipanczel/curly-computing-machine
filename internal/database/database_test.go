package database

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func mustStartMongoContainer() (func(context.Context) error, error) {
	dbContainer, err := mongodb.Run(context.Background(), "mongo:latest")
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	database = "curly_test"

	teardown, err := mustStartMongoContainer()
	if err != nil {
		log.Fatalf("could not start mongodb container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mongodb container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func (s *service) deleteAuthorColl(ctx context.Context) error {
	filter := bson.M{}

	_, err := s.authorsColl.DeleteMany(ctx, filter)
	return err
}

func (s *service) deleteBookColl(ctx context.Context) error {
	filter := bson.M{}

	_, err := s.booksColl.DeleteMany(ctx, filter)
	return err
}

func (s *service) deleteBorrowerColl(ctx context.Context) error {
	filter := bson.M{}

	_, err := s.borrowersColl.DeleteMany(ctx, filter)
	return err
}