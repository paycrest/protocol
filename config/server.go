package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServerConfiguration type defines the server configurations
type ServerConfiguration struct {
	Debug                    bool
	Host                     string
	Port                     string
	Timezone                 string
	AllowedHosts             string
	Environment              string
	SentryDSN                string
	HostDomain               string
	RateLimitUnauthenticated int
	RateLimitAuthenticated   int
	QuidaxURL                string
	BitgetURL                string
	BinanceURL               string
}

// ServerConfig sets the server configuration
func ServerConfig() *ServerConfiguration {
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("SERVER_TIMEZONE", "Africa/Lagos")
	viper.SetDefault("ALLOWED_HOSTS", "*")
	viper.SetDefault("ENVIRONMENT", "local")
	viper.SetDefault("SENTRY_DSN", "")
	viper.SetDefault("RATE_LIMIT_UNAUTHENTICATED", 5)
	viper.SetDefault("RATE_LIMIT_AUTHENTICATED", 50)
	viper.SetDefault("QUIDAX_URL", "https://www.quidax.com/api/v1/markets")
	viper.SetDefault("BITGET_URL", "https://api.bitget.com/api/mix/v1/market/p2p/advertisements")
	viper.SetDefault("BINANCE_URL", "https://api.binance.com/api/v3/ticker/price")

	return &ServerConfiguration{
		Debug:                    viper.GetBool("DEBUG"),
		Host:                     viper.GetString("SERVER_HOST"),
		Port:                     viper.GetString("SERVER_PORT"),
		Timezone:                 viper.GetString("SERVER_TIMEZONE"),
		AllowedHosts:             viper.GetString("ALLOWED_HOSTS"),
		Environment:              viper.GetString("ENVIRONMENT"),
		SentryDSN:                viper.GetString("SENTRY_DSN"),
		HostDomain:               viper.GetString("HOST_DOMAIN"),
		RateLimitUnauthenticated: viper.GetInt("RATE_LIMIT_UNAUTHENTICATED"),
		RateLimitAuthenticated:   viper.GetInt("RATE_LIMIT_AUTHENTICATED"),
		QuidaxURL:                viper.GetString("QUIDAX_URL"),
		BitgetURL:                viper.GetString("BITGET_URL"),
		BinanceURL:               viper.GetString("BINANCE_URL"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		panic(fmt.Sprintf("config SetupConfig() error: %s", err))
	}
}
