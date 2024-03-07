package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// EmailConfiguration defines the email service configurations
type NotificationConfiguration struct {
	EmailDomain         string
	EmailAPIKey         string
}

// EmailConfig sets the email configurations
func NotificationConfig() (config *NotificationConfiguration) {
	viper.SetDefault("EMAIL_DOMAIN", "sandbox9c66b379b78d43d2b1533bf2a09a5325.mailgun.org")

	return &NotificationConfiguration{
		EmailDomain:         viper.GetString("EMAIL_DOMAIN"),
		EmailAPIKey:         viper.GetString("EMAIL_API_KEY"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		panic(fmt.Sprintf("config SetupConfig() error: %s", err))
	}
}
