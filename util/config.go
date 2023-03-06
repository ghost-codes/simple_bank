package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SecretKey           string        `mapstructure:"SECRET_KEY"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddr          string        `mapstructure:"SEVER_ADDR"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	fmt.Println(config)
	return config, err
}
