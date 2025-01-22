package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/routers"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/tasks"
	"github.com/paycrest/protocol/utils/logger"
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

	// err := tasks.FixDatabaseMisHap()
	// if err != nil {
	// 	logger.Errorf("FixDatabaseMisHap: %v", err)
	// }

	// Initialize Redis
	if err := storage.InitializeRedis(); err != nil {
		log.Println(err)
		logger.Fatalf("Redis initialization: %v", err)
	}

	// Subscribe to Redis keyspace events
	tasks.SubscribeToRedisKeyspaceEvents()

	// Start cron jobs
	tasks.StartCronJobs()

	// Run the server
	router := routers.Routes()

	appServer := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof("Server Running at :%v", appServer)

	logger.Fatalf("%v", router.Run(appServer))
}
