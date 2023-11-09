package config

import (
	"time"

	"github.com/paycrest/protocol/utils/logger"
	"github.com/spf13/viper"
)

// AuthConfiguration defines the authentication & authorization settings
type AuthConfiguration struct {
	Secret                 string
	JwtAccessHourLifespan  time.Duration
	JwtRefreshHourLifespan time.Duration
	HmacTimestampAge       time.Duration
	PasswordResetLifespan  time.Duration
}

// AuthConfig sets the authentication & authorization configurations
func AuthConfig() (config *AuthConfiguration) {
	viper.SetDefault("JWT_ACCESS_HOUR_LIFESPAN", 15)
	viper.SetDefault("JWT_REFRESH_HOUR_LIFESPAN", 24)
	viper.SetDefault("HMAC_TIMESTAMP_AGE", 5)
	viper.SetDefault("PASSWORD_RESET_LIFESPAN", 5)

	return &AuthConfiguration{
		Secret:                 viper.GetString("SECRET"),
		JwtAccessHourLifespan:  time.Duration(viper.GetInt("JWT_ACCESS_HOUR_LIFESPAN")) * time.Minute,
		JwtRefreshHourLifespan: time.Duration(viper.GetInt("JWT_REFRESH_HOUR_LIFESPAN")) * time.Hour,
		HmacTimestampAge:       time.Duration(viper.GetInt("HMAC_TIMESTAMP_AGE")) * time.Minute,
		PasswordResetLifespan:  time.Duration(viper.GetInt("PASSWORD_RESET_LIFESPAN")) * time.Minute,
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
