package main

import (
	"github.com/akmamun/gin-boilerplate-examples/config"
	"github.com/akmamun/gin-boilerplate-examples/infra/database"
	"github.com/akmamun/gin-boilerplate-examples/infra/logger"
	"github.com/akmamun/gin-boilerplate-examples/routers"
	"github.com/spf13/viper"
	"time"
)

func main() {
	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DBConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	router := routers.Routes()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
