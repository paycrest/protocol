package config

import (
	"os"

	"github.com/cosmos/go-bip39"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// Configuration type
type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Auth     AuthConfiguration
	Order    OrderConfiguration
	Notification    NotificationConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.AddConfigPath("../..")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")

	envFilePath := os.Getenv("ENV_FILE_PATH")
	if envFilePath == "" {
		envFilePath = ".env" // Set default value to ".env"
	}

	viper.SetConfigName(envFilePath)
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

	var serverConf = ServerConfig()

	valid := bip39.IsMnemonicValid(serverConf.HDWalletMnemonic)
	if !valid {
		logger.Errorf("Invalid mnemonic phrase")
		return nil
	}

	return nil
}
