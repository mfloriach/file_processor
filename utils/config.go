package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Compress Fields `mapstructure:"compress"`
	Metadata Fields `mapstructure:"metadata"`
	Reader   Fields `mapstructure:"reader"`
	Store    Fields `mapstructure:"store"`
	Injestor Fields `mapstructure:"injestor"`
}

type Fields struct {
	QueueLen   int `mapstructure:"queue_length"`
	WorkersNum int `mapstructure:"num_workers"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}

	x := viper.New()
	x.SetConfigType("json")
	x.SetConfigFile("config.json")

	if err := x.ReadInConfig(); err != nil {
		log.Fatal("Can't find the file config.json: ", err)
	}

	if err := x.Unmarshal(&cfg); err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return cfg
}
