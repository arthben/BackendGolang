package config

import (
	"os"

	"github.com/spf13/viper"
)

type EnvParams struct {
	App struct {
		Port        string `yaml:"port"`
		WaitTimeOut int    `yaml:"waitTimeout"`
	} `yaml:"app"`
	DB struct {
		DSN     string `yaml:"dsn"`
		MinPool int    `yaml:"minPool"`
		MaxPool int    `yaml:"maxPool"`
	} `yaml:"db"`
	OpenWeather struct {
		APIKey string `yaml:"apikey"`
	} `yaml:"openweather"`
}

func LoadConfig() (*EnvParams, error) {
	// read config file
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// read OS env to fill config file parameter
	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		viper.Set(key, os.ExpandEnv(val))
	}

	var cfg EnvParams
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
