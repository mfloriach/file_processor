package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBConn             string `mapstructure:"DB_CONN"`
	DBName             string `mapstructure:"DB_NAME"`
	MinioEndpoint      string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey     string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey     string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucketNameSrc string `mapstructure:"MINIO_BUCKET_NAME_SRC"`
	MinioBucketNameOut string `mapstructure:"MINIO_BUCKET_NAME_OUT"`
	MinioLocation      string `mapstructure:"MINIO_LOCATION"`
	ImdbToken          string `mapstructure:"IMDB_TOKEN"`
	ImdbBaseUrl        string `mapstructure:"IMDB_BASE_URL"`
}

var env *Env

func GetEnv() *Env {
	if env != nil {
		return env
	}

	x := viper.New()
	x.SetConfigFile(".env")

	if err := x.ReadInConfig(); err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	if err := x.Unmarshal(&env); err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return env
}
