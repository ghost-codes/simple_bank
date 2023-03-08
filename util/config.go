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
	DBHOST              string        `mapstructure:"DB_HOST"`
	DBPort              string        `mapstructure:"DB_PORT"`
	DBUser              string        `mapstructure:"DB_USER"`
	DBPassword          string        `mapstructure:"DB_PASSWORD"`
	DBName              string        `mapstructure:"DB_NAME"`
}

// func (config *Config) DBSource() string {
// 	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", config.DBUser, config.DBPassword, config.DBHOST, config.DBPort, config.DBName)
// }

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
