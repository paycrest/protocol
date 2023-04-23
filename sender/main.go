package main

import (
	"time"

	"github.com/paycrest/paycrest-services/sender/config"
	"github.com/paycrest/paycrest-services/sender/database"
	"github.com/paycrest/paycrest-services/sender/routers"
	"github.com/paycrest/paycrest-services/sender/utils/logger"

	"github.com/spf13/viper"
)

func main() {
	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	DSN := config.DbConfiguration()

	if err := database.DBConnection(DSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	router := routers.Routes()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
