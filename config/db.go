package config

import (
	"fmt"

	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// DatabaseConfiguration type defines the server configurations
type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

// DBConfiguration sets the database configuration
func DBConfiguration() (DSN string) {
	DbName := viper.GetString("DB_NAME")
	DbUser := viper.GetString("DB_USER")
	DbPassword := viper.GetString("DB_PASSWORD")
	DbHost := viper.GetString("DB_HOST")
	DbPort := viper.GetString("DB_PORT")
	DbSslMode := viper.GetString("SSL_MODE")

	DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DbHost, DbUser, DbPassword, DbName, DbPort, DbSslMode,
	)

	return
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
