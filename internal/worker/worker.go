package worker

import (
	"context"
	"strings"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/monitors"
	"github.com/brunohfonseca/ratatoskr/internal/notifications"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// StartWorker inicia o worker que consome endpoints do Redis Stream
func StartWorker(ctx context.Context, redisClient *redis.Client, groupName, consumerName string) {
	streams := []string{"notifications", "endpoints", "ssl-checks"}

	log.Info().
		Str("group", groupName).
		Str("consumer", consumerName).
		Msg("🚀 Starting")

	// Cria os grupos (ou confirma que já existem)
	for _, s := range streams {
		if err := redisClient.XGroupCreateMkStream(ctx, s, groupName, "0").Err(); err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
			log.Fatal().Err(err).Str("stream", s).Msg("❌ Failed to create consumer group")
		}
	}

	// [keys..., ids...]
	streamArgs := append([]string{}, streams...)
	for range streams {
		streamArgs = append(streamArgs, ">")
	}

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("🛑 Worker shutting down")
			return

		default:
			results, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  streamArgs,
				Block:    time.Second,
				Count:    10,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue // timeout sem mensagens
				}
				log.Fatal().Err(err).Msg("❌ Fatal error reading from stream")
			}

			for _, res := range results {
				for _, msg := range res.Messages {
					log.Info().
						Str("stream", res.Stream).
						Str("id", msg.ID).
						Msg("📨 Mensagem recebida")

					switch res.Stream {
					case "endpoints":
						monitors.ProcessEndpoint(ctx, redisClient, res.Stream, groupName, msg)
					case "ssl-checks":
						monitors.ProcessSSLCheck(ctx, redisClient, res.Stream, groupName, msg)
					case "notifications":
						notifications.ProcessNotification(ctx, redisClient, res.Stream, groupName, msg)
					}

					if _, err := redisClient.XAck(ctx, res.Stream, groupName, msg.ID).Result(); err != nil {
						log.Fatal().Err(err).Str("id", msg.ID).Msg("❌ Failed to ACK message")
					}
					if _, err := redisClient.XDel(ctx, res.Stream, msg.ID).Result(); err != nil {
						log.Fatal().Err(err).Str("id", msg.ID).Msg("❌ Failed to delete message")
					}
				}
			}
		}
	}
}
