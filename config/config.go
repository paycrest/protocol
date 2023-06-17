package config

import (
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// Configuration type
type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Auth     AuthConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.AddConfigPath("../..")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	return nil
}
