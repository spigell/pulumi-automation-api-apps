package manager

import (
	"errors"

	"github.com/spf13/viper"
)

var ErrWrongMaxKeep = errors.New("max-keep must be greater than 0")

type Config struct {
	Verbose       bool
	APIServerPort int `mapstructure:"api-server-port"`
	Stack         Stack
	Token         string
	MaxKeep       int `mapstructure:"max-keep"`
	Diff          bool
}

type Stack struct {
	Name string
	Path string
}

func GetConfig() (*Config, error) {
	var config *Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	config.Token = viper.GetString("hcloud-token")
	if config.Token == "" {
		config.Token = viper.GetString("hcloud_token")
	}

	if config.MaxKeep < 1 {
		return nil, ErrWrongMaxKeep
	}

	return config, nil
}
