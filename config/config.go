package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {

	config = viper.New()
	config.SetConfigName("config")    // name of config file (without extension)
	config.SetConfigType("yml")       // REQUIRED if the config file does not have the extension in the name
	config.AddConfigPath("./config")  // optionally look for config in the working directory
	config.AddConfigPath(".")         // optionally look for config in the working directory
	config.AddConfigPath("../config") // optionally look at parent dir then config
	config.AutomaticEnv()             // auto read env variables

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func GetConfig() *viper.Viper {
	return config
}
