package worker

import (
	"context"
	"strings"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/monitors"
	"github.com/brunohfonseca/ratatoskr/internal/notifications"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
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
			logger.FatalLog("❌ Failed to create consumer group", err)
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
			logger.InfoLog("🛑 Worker shutting down")
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
				logger.FatalLog("❌ Fatal error reading from stream", err)
			}

			for _, res := range results {
				for _, msg := range res.Messages {
					switch res.Stream {
					case "endpoints":
						monitors.ProcessEndpoint(msg)
						logger.DebugLog("✅ Endpoint processed")
					case "ssl-checks":
						monitors.ProcessSSLCheck(msg)
						logger.DebugLog("✅ SSL check processed")
					case "notifications":
						notifications.ProcessNotification(msg)
						logger.DebugLog("✅ Notification processed")
					}

					if _, err := redisClient.XAck(ctx, res.Stream, groupName, msg.ID).Result(); err != nil {
						logger.FatalStrLog("❌ Failed to ACK message", "id", msg.ID)
					}
					if _, err := redisClient.XDel(ctx, res.Stream, msg.ID).Result(); err != nil {
						logger.FatalStrLog("❌ Failed to delete message", "id", msg.ID)
					}
				}
			}
		}
	}
}
