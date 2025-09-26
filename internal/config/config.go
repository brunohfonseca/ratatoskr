package config

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var appConfig *AppConfig

func SetupLogs() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 -0700",
	})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

type AppConfig struct {
	Server struct {
		Port int `yaml:"port"`
		SSL  struct {
			Enabled bool   `yaml:"enabled"`
			Port    int    `yaml:"port"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"ssl"`
	} `yaml:"server"`
	Database struct {
		MongoURL string `yaml:"mongo_url"`
	} `yaml:"database"`
	Redis struct {
		RedisURL string `yaml:"redis_url"`
	} `yaml:"redis"`
	Alerts struct {
		Slack struct {
			Channel string `yaml:"channel"`
			Token   string `yaml:"token"`
		} `yaml:"slack"`
		Telegram struct {
			BotToken string `yaml:"bot_token"`
			ChatID   string `yaml:"chat_id"`
		} `yaml:"telegram"`
	} `yaml:"alerting"`
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Database.MongoURL == "" || cfg.Redis.RedisURL == "" {
		return nil, errors.New("missing required variables in config")
	}

	appConfig = &cfg
	return appConfig, nil
}

func Get() *AppConfig {
	return appConfig
}
