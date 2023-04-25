package main

import (
	"time"

	"github.com/paycrest/paycrest-protocol/sender/config"
	"github.com/paycrest/paycrest-protocol/sender/database"
	"github.com/paycrest/paycrest-protocol/sender/routers"
	"github.com/paycrest/paycrest-protocol/sender/utils/logger"

	"github.com/spf13/viper"
)

func main() {
	// Set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	// Setup config
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	// Connect to the database
	DSN := config.DBConfiguration()

	if err := database.DBConnection(DSN); err != nil {
		logger.Fatalf("database DBConnection error: %s", err)
	}

	defer database.GetClient().Close()

	// Run the server
	router := routers.Routes()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
