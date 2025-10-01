package worker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// StartWorker inicia o worker que consome endpoints do Redis Stream
func StartWorker(ctx context.Context, redisClient *redis.Client, groupName, consumerName string) {
	streams := []string{"alerts", "endpoints", "ssl-checks"}

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
						processEndpoint(ctx, redisClient, res.Stream, groupName, msg)
					case "ssl-checks":
						// processSSLCheck(ctx, redisClient, res.Stream, groupName, msg)
					}

					if _, err := redisClient.XAck(ctx, res.Stream, groupName, msg.ID).Result(); err != nil {
						log.Fatal().Err(err).Str("id", msg.ID).Msg("❌ Failed to ACK message")
					}
				}
			}
		}
	}
}

// processEndpoint processa um endpoint do stream
func processEndpoint(ctx context.Context, redisClient *redis.Client, stream, group string, msg redis.XMessage) {
	uuid := msg.Values["uuid"].(string)
	domain := msg.Values["domain"].(string)
	path, _ := msg.Values["path"].(string)
	checkSSLStr, _ := msg.Values["check_ssl"].(string)
	checkSSL := checkSSLStr == "true"

	log.Info().
		Str("uuid", uuid).
		Str("domain", domain).
		Msg("🔍 Processing endpoint")

	// Faz health check
	status, responseTime := doHealthCheck(domain, path, checkSSL)

	log.Info().
		Str("uuid", uuid).
		Str("status", status).
		Int64("response_time_ms", responseTime).
		Msg("✅ Health check completed")

	// TODO: Salvar resultado no banco de dados
	//db.Exec("UPDATE endpoints SET status = $1, last_check = NOW() WHERE uuid = $2", status, uuid)

	// ACK da mensagem (marca como processada)
	redisClient.XAck(ctx, stream, group, msg.ID)

	// Remove a mensagem do stream (não precisa manter histórico no Redis)
	redisClient.XDel(ctx, stream, msg.ID)
}

// doHealthCheck faz uma requisição HTTP e retorna status e tempo de resposta
func doHealthCheck(domain, path string, checkSSL bool) (string, int64) {
	protocol := "http"
	if checkSSL {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s%s", protocol, domain, path)
	if path == "" {
		url = fmt.Sprintf("%s://%s", protocol, domain)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	responseTime := time.Since(start).Milliseconds()

	if err != nil {
		log.Warn().Err(err).Str("url", url).Msg("Health check failed")
		return "offline", responseTime
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "online", responseTime
	}

	return "offline", responseTime
}
