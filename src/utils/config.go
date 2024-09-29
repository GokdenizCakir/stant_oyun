package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	PORT              string `mapstructure:"PORT"`
	JWT_SECRET_ACCESS string `mapstructure:"JWT_SECRET_ACCESS"`
	QUESTION_COUNT    string `mapstructure:"QUESTION_COUNT"`
	QUESTION_SECONDS  string `mapstructure:"QUESTION_SECONDS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err

	}

	err = viper.Unmarshal(&config)
	return
}
