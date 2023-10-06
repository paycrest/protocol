package main

import (
	"fmt"
	"time"

	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/routers"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/tasks"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

func main() {
	// Set timezone
	conf := config.ServerConfig()
	loc, _ := time.LoadLocation(conf.Timezone)
	time.Local = loc

	// Connect to the database
	DSN := config.DBConfig()

	if err := storage.DBConnection(DSN); err != nil {
		logger.Fatalf("database DBConnection: %s", err)
	}

	defer storage.GetClient().Close()

	// Seed the database
	if err := storage.SeedAll(); err != nil {
		logger.Fatalf("database SeedAll: %s", err)
	}

	// Initialize Redis
	if err := storage.InitializeRedis(); err != nil {
		logger.Fatalf("Redis initialization: %s", err)
	}

	// Start indexer
	if err := tasks.ContinueIndexing(); err != nil {
		logger.Fatalf("ContinueIndexing: %s", err)
	}

	// Start processing orders
	if err := tasks.ProcessOrders(); err != nil {
		logger.Fatalf("ProcessOrders: %s", err)
	}

	// Start scanning unfulfilled lockOrders & reassignment.
	if err := tasks.ProcessUnfulfilledLockOrders(); err != nil {
		logger.Fatalf("UnfulfilledReassignment: %s", err)
	}

	// Subscribe to Redis keyspace events
	tasks.SubscribeToRedisKeyspaceEvents()

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :%v", appServer)

	logger.Fatalf("%v", router.Run(appServer))
}
