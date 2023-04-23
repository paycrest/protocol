package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func DbConfiguration() (DSN string) {
	DBName := viper.GetString("DB_NAME")
	DBUser := viper.GetString("DB_USER")
	DBPassword := viper.GetString("DB_PASSWORD")
	DBHost := viper.GetString("DB_HOST")
	DBPort := viper.GetString("DB_PORT")
	DBSslMode := viper.GetString("SSL_MODE")

	DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		DBHost, DBUser, DBPassword, DBName, DBPort, DBSslMode,
	)

	return
}
