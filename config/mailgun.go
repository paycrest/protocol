package config

import (
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// MailGunConfiguration defines the mailgun configurations
type MailGunConfiguration struct {
	Domain string
	ApiKey string
}

// MailGunConfig sets the mailgun configurations
func MailGunConfig() (config *MailGunConfiguration) {
	viper.SetDefault("MAILGUN_DOMAIN", "sandbox9c66b379b78d43d2b1533bf2a09a5325.mailgun.org")

	return &MailGunConfiguration{
		Domain: viper.GetString("MAILGUN_DOMAIN"),
		ApiKey: viper.GetString("MAILGUN_APIKEY"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
