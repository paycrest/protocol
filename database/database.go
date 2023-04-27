package database

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/migrate"

	_ "github.com/lib/pq"
)

var (
	Client *ent.Client
	Err    error
)

// DBConnection create database connection
func DBConnection(DSN string) error {
	var client = Client

	client, err := ent.Open(dialect.Postgres, DSN)
	if err != nil {
		Err = err
		log.Println("Database connection error")
		return err
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		return err
	}

	Client = client

	return nil
}

// GetClient connection
func GetClient() *ent.Client {
	return Client
}

// GetError connection error
func GetError() error {
	return Err
}
