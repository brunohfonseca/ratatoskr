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

	log.Info().
		Str("group", groupName).
		Str("consumer", consumerName).
		Msg("Initializing worker...")

	// Cria consumer groups para cada stream
	streams := []string{"alerts", "endpoints", "ssl-checks"}

	for _, stream := range streams {
		err := redisClient.XGroupCreateMkStream(ctx, stream, groupName, "0").Err()
		if err != nil {
			// Ignora se grupo já existe
			if err.Error() == "BUSYGROUP Consumer Group name already exists" {
				log.Info().Str("stream", stream).Str("group", groupName).Msg("Consumer group already exists")
			} else {
				log.Error().Err(err).Str("stream", stream).Str("group", groupName).Msg("Failed to create consumer group")
			}
		} else {
			log.Info().Str("stream", stream).Str("group", groupName).Msg("Consumer group created")
		}
	}

	// Pequeno delay para garantir que Redis processou
	time.Sleep(100 * time.Millisecond)

	log.Info().Str("consumer", consumerName).Str("group", groupName).Msg("🚀 Consumer started")

	for {
		// Lê mensagens do stream
		results, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{"alerts", "endpoints", "ssl-checks", ">", ">", ">"},
			Count:    10,
			Block:    1 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				// Se for erro de grupo não encontrado, loga e para
				if err.Error() == "NOGROUP No such key '>' or consumer group '"+groupName+"' in XREADGROUP with GROUP option" {
					log.Fatal().
						Err(err).
						Str("group", groupName).
						Str("consumer", consumerName).
						Msg("❌ Consumer group not found - stopping worker")
				}

				// Outros erros também param
				log.Fatal().
					Err(err).
					Str("group", groupName).
					Str("consumer", consumerName).
					Msg("❌ Fatal error reading from stream - stopping worker")
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
					processEndpoint(ctx, redisClient, streamName, groupName, msg)
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
