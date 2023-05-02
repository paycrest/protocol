package config

import (
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Debug        bool
	Host         string
	Port         string
	Timezone     string
	AllowedHosts string
}

func ServerConfig() *ServerConfiguration {
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("SERVER_TIMEZONE", "Africa/Lagos")
	viper.SetDefault("ALLOWED_HOSTS", "*")

	return &ServerConfiguration{
		Debug:        viper.GetBool("DEBUG"),
		Host:         viper.GetString("SERVER_HOST"),
		Port:         viper.GetString("SERVER_PORT"),
		Timezone:     viper.GetString("SERVER_TIMEZONE"),
		AllowedHosts: viper.GetString("ALLOWED_HOSTS"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
