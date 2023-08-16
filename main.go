package main

import (
	"fmt"
	"time"

	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/routers"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

func main() {
	// Set timezone
	conf := config.ServerConfig()
	loc, _ := time.LoadLocation(conf.Timezone)
	time.Local = loc

	// Connect to the database
	DSN := config.DBConfiguration()

	if err := database.DBConnection(DSN); err != nil {
		logger.Fatalf("database DBConnection error: %s", err)
	}

	defer database.GetClient().Close()

	// Seed the database
	if err := database.SeedAll(); err != nil {
		logger.Fatalf("database SeedAll error: %s", err)
	}

	// Start indexer
	if err := ContinueIndexing(); err != nil {
		logger.Fatalf("continueIndexing error: %s", err)
	}

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :%v", appServer)

	logger.Fatalf("%v", router.Run(appServer))
}
