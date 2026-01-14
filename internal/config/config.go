package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	DBURL          string `mapstructure:"DATABASE_URL"`
	RedisURL       string `mapstructure:"REDIS_URL"`
	GithubClientID string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubSecret   string `mapstructure:"GITHUB_CLIENT_SECRET"`
	EncryptionKey  string `mapstructure:"ENCRYPTION_KEY"`
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
