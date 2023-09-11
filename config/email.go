package config

import (
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// EmailConfiguration defines the email service configurations
type EmailConfiguration struct {
	Domain string
	ApiKey string
}

// EmailConfig sets the email configurations
func EmailConfig() (config *EmailConfiguration) {
	viper.SetDefault("EMAIL_DOMAIN", "sandbox9c66b379b78d43d2b1533bf2a09a5325.mailgun.org")

	return &EmailConfiguration{
		Domain: viper.GetString("EMAIL_DOMAIN"),
		ApiKey: viper.GetString("EMAIL_APIKEY"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
