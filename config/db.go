package config

import (
	"fmt"

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

// DBConfig sets the database configuration
func DBConfig() (DSN string) {
	DbName := viper.GetString("DB_NAME")
	DbUser := viper.GetString("DB_USER")
	DbPassword := viper.GetString("DB_PASSWORD")
	DbHost := viper.GetString("DB_HOST")
	DbPort := viper.GetString("DB_PORT")
	DbSslMode := viper.GetString("SSL_MODE")

	DSN = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", DbUser, DbPassword, DbHost, DbPort, DbName, DbSslMode)

	return
}

func init() {
	if err := SetupConfig(); err != nil {
		panic(fmt.Sprintf("config SetupConfig() error: %s", err))
	}
}
