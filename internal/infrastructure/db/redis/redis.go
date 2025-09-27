package infra

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var RedisClient *redis.Client

func ConnectRedis(RedisURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts, err := redis.ParseURL(RedisURL)
	if err != nil {
		log.Fatal().Msgf("Error parsing Redis URL: %s", err)
	}

	client := redis.NewClient(opts)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Msgf("Error connecting to Redis: %s", err)
	}
	log.Info().Msg("✅ Connected to Redis")
	RedisClient = client
}

func DisconnectRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Fatal().Msgf("Erro ao desconectar do Redis: %v", err)
		} else {
			log.Info().Msg("✅ Disconnected from Redis")
		}
	}
}

// CheckRedisHealth verifica o status da conexão com Redis
func CheckRedisHealth() (bool, string, error) {
	if RedisClient == nil {
		return false, "disconnected", nil // Não é um erro fatal, apenas não conectado
	}

	// Criar contexto com timeout curto para health check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Fazer ping para verificar se a conexão está ativa
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return false, "error", err
	}

	return true, "connected", nil
}
