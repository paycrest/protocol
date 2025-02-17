package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// StreamConfiguration type defines the stream configurations

type StreamConfiguration struct {
	QuickNodeAPIKey 			string
	QuickNodeStreamAPIURL 		string
	QuickNodeFunctionAPIURL 	string
	QuickNodePrivateKey 		string
}

// ServerConfig sets the server configuration
func StreamConfig() *StreamConfiguration {
	viper.SetDefault("QUICKNODE_API_KEY", "")
	viper.SetDefault("QUICKNODE_STREAM_API_URL", "")
	viper.SetDefault("QUICKNODE_FUNCTION_API_URL", "")
	viper.SetDefault("QUICKNODE_PRIVATE_KEY", "")

	return &StreamConfiguration{
		QuickNodeAPIKey: 			viper.GetString("QUICKNODE_API_KEY"),
		QuickNodeStreamAPIURL: 		viper.GetString("QUICKNODE_STREAM_API_URL"),
		QuickNodeFunctionAPIURL: 	viper.GetString("QUICKNODE_FUNCTION_API_URL"),
		QuickNodePrivateKey: 		viper.GetString("QUICKNODE_PRIVATE_KEY"),
	}
}

func init() {
	if err := SetupConfig(); err != nil {
		panic(fmt.Sprintf("config SetupConfig() error: %s", err))
	}
}
