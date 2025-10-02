package config

import (
	"errors"
	"os"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
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
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"` // development, production
	Server      struct {
		Port           int      `yaml:"port"`
		TrustedProxies []string `yaml:"trusted_proxies"`
		SSL            struct {
			Enabled bool   `yaml:"enabled"`
			Port    int    `yaml:"port"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"ssl"`
	} `yaml:"server"`
	Database struct {
		PostgresURL string `yaml:"postgres_url"`
	} `yaml:"database"`
	Redis struct {
		RedisURL string `yaml:"redis_url"`
	} `yaml:"redis"`
	JWT struct {
		JWTSecret          string `yaml:"jwt_secret"`
		JWTExpirationHours int    `yaml:"jwt_expiration_hours"`
	} `yaml:"jwt"`
	OIDC struct {
		Enabled      bool   `yaml:"enabled"`
		URL          string `yaml:"url"`
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectURL  string `yaml:"redirect_url"`
	} `yaml:"oidc"`
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

func LoadYamlConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Database.PostgresURL == "" || cfg.Redis.RedisURL == "" {
		return nil, errors.New("missing required variables in config")
	}

	appConfig = &cfg
	return appConfig, nil
}

func Get() *AppConfig {
	return appConfig
}

func LoadConfig(path string) (string, error) {
	_, err := LoadYamlConfig(path)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("🚀 Iniciando o serviço com o arquivo de configuração: %s", path)
	redis.ConnectRedis(appConfig.Redis.RedisURL)
	postgres.ConnectPostgres(appConfig.Database.PostgresURL)
	return "", nil
}
