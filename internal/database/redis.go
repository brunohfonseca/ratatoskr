package database

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
