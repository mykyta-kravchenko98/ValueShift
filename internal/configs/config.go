package configs

import (
	"github.com/spf13/viper"
)

type ExchangeApiConfig struct {
	ApiKey string `mapstructure:"api-key"`
	URL    string `mapstructure:"url"`
}

type MongoDBConfig struct {
	URL string `mapstructure:"url"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	ExchangeApi ExchangeApiConfig `mapstructure:"exchange-api"`
	MongoDB     MongoDBConfig     `mapstructure:"mongo-db"`
	Server      ServerConfig      `mapstructure:"server"`
}

var (
	vp     *viper.Viper
	config *Config
)

func LoadConfigs(env string) (*Config, error) {
	vp = viper.New()

	vp.SetConfigType("json")
	vp.SetConfigName(env)
	vp.AddConfigPath("../configs/")
	vp.AddConfigPath("../../configs/")
	vp.AddConfigPath("configs/")

	err := vp.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return &Config{}, err
	}

	return config, err
}

func GetConfig() *Config {
	return config
}
