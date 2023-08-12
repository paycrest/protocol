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

	// TODO: Continue indexing erc20 transfers which could have been impacted by a downtime
	// we can find them by checking receive addresses with a payment order that are
	// still unused or partial after the receive address validity period.

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :%v", appServer)

	logger.Fatalf("%v", router.Run(appServer))
}
