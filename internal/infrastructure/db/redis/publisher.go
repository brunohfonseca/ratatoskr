package infra

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func PublishAlert() {}

func StreamPublish(ctx context.Context, msg *redis.XAddArgs) error {
	redisClient := RedisClient
	_, err := redisClient.XAdd(ctx, msg).Result()
	if err != nil {
		log.Error().Err(err).Msg("⚠️ Erro ao publicar no Redis")
	}
	return err
}
