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
	// Setup config
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

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

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :", appServer)

	logger.Fatalf("%v", router.Run(appServer))

}
