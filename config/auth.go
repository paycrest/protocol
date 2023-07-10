package config

import (
	"time"

	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/spf13/viper"
)

// AuthConfiguration defines the authentication & authorization settings
type AuthConfiguration struct {
	Secret                 string
	JwtAccessHourLifespan  time.Duration
	JwtRefreshHourLifespan time.Duration
	HmacTimestampAge       int64
}

// AuthConfig sets the authentication & authorization configurations
func AuthConfig() (config *AuthConfiguration) {
	viper.SetDefault("JWT_ACCESS_HOUR_LIFESPAN", 15)
	viper.SetDefault("JWT_REFRESH_HOUR_LIFESPAN", 24)

	return &AuthConfiguration{
		Secret:                 viper.GetString("SECRET"),
		JwtAccessHourLifespan:  time.Duration(viper.GetInt("JWT_ACCESS_HOUR_LIFESPAN")) * time.Minute,
		JwtRefreshHourLifespan: time.Duration(viper.GetInt("JWT_REFRESH_HOUR_LIFESPAN")) * time.Hour,
		HmacTimestampAge:       5,
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
}
