package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	CRConfig ChatRoomConfig
)

type Http struct {
	Address string
	Port    int
}

type ChatRoomConfig struct {
	Http Http
}

func InitConfig(configFile string) {
	viper.SetConfigName("configFile")
	viper.AddConfigPath(configFile)

	viper.SetDefault("Http.Port", 2023)
	// viper.SetDefault("Http.Address", "localhost")

	if err := viper.ReadInConfig(); err != nil {
		// panic(fmt.Errorf("Fatal error read config file: %s \n", err))
	}

	if err := viper.Unmarshal(&CRConfig); err != nil {
		// panic(fmt.Errorf("Fatal error unmarshal config file: %s \n", err))
	}

	fmt.Printf("Http config: %+v\n", CRConfig.Http)
}
