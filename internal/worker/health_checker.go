package worker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// StartHealthCheckWorker inicia o worker que consome endpoints do Redis Stream
func StartHealthCheckWorker(redisClient *redis.Client, groupName, consumerName string) {
	ctx := context.Background()
	group := groupName
	consumer := consumerName

	redisClient.XGroupCreateMkStream(ctx, "alerts", group, "0")
	redisClient.XGroupCreateMkStream(ctx, "endpoints", group, "0")
	redisClient.XGroupCreateMkStream(ctx, "ssl-checks", group, "0")

	log.Info().Msg("🚀 Health Check Worker started")

	for {
		// Lê mensagens do stream
		results, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{"alerts", ">", "endpoints", ">", "ssl-checks", ">"},
			Count:    10,
			Block:    1 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				log.Error().Err(err).Msg("Failed to read from stream")
			}
			continue
		}

		// Processa mensagens
		for _, result := range results {
			streamName := result.Stream // "endpoints" ou "ssl-checks"

			for _, msg := range result.Messages {
				log.Info().Str("stream", streamName).Msg("📨 Mensagem recebida")

				// Processa baseado em qual stream veio
				if streamName == "endpoints" {
					processEndpoint(ctx, redisClient, streamName, group, msg)
				} else if streamName == "ssl-checks" {
					//processSSLCheck(ctx, redisClient, streamName, group, msg)
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
