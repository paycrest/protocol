package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/migrate"
	_ "github.com/paycrest/protocol/ent/runtime" // ent runtime

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	// Client holds the database connection
	Client *ent.Client
	// Err holds database connection error
	Err error
)

// DBConnection create database connection
func DBConnection(DSN string) error {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		Err = err
		log.Println("Database connection error")
		return err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)

	// Integrate sql.DB to ent.Client.
	client := ent.NewClient(ent.Driver(drv))

	conf := config.ServerConfig()

	// Run the auto migration tool.
	if conf.Environment == "local" {
		if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
			return err
		}
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
