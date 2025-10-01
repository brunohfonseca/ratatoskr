package monitors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func ProcessEndpoint(ctx context.Context, redisClient *redis.Client, stream, group string, msg redis.XMessage) {
	var endpoint models.Endpoint

	// Leitura segura dos campos do Redis
	uuid, _ := msg.Values["uuid"].(string)
	domain, _ := msg.Values["domain"].(string)
	path, _ := msg.Values["path"].(string)
	checkSSLStr, _ := msg.Values["check_ssl"].(string)

	// Timeout com valor padrão de 30 segundos
	timeout := 30
	if timeoutVal, ok := msg.Values["timeout"].(int64); ok {
		timeout = int(timeoutVal)
	} else if timeoutVal, ok := msg.Values["timeout"].(int); ok {
		timeout = timeoutVal
	}

	url := fmt.Sprintf("%s/%s", domain, path)
	doHealthCheck(url, timeout)

	services.UpdateCheck(endpoint)

	if checkSSLStr == "true" {
		_, err := FetchSSL(domain)
		if err != nil {
			return
		}
	}

	log.Info().Msgf("✅ Health check completed in %s", uuid)
}

func doHealthCheck(url string, timeout int) models.EndpointResponse {
	start := time.Now()

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return models.EndpointResponse{
			ResponseStatusCode: 0, // Sem resposta devido ao erro
			ResponseMessage:    err.Error(),
			ResponseTime:       int(time.Since(start).Milliseconds()),
		}
	}
	defer resp.Body.Close()

	return models.EndpointResponse{
		ResponseStatusCode: resp.StatusCode,
		ResponseMessage:    "Success",
		ResponseTime:       int(time.Since(start).Milliseconds()),
	}
}
