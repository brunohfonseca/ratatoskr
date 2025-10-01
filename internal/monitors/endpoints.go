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
	uuid := msg.Values["uuid"].(string)
	domain := msg.Values["domain"].(string)
	path := msg.Values["path"].(string)
	timeout := msg.Values["timeout"].(int)
	checkSSLStr, _ := msg.Values["check_ssl"].(string)

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
		Timeout: time.Duration(timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return models.EndpointResponse{
			ResponseStatusCode: resp.StatusCode,
			ResponseMessage:    err.Error(),
		}
	}
	defer resp.Body.Close()

	return models.EndpointResponse{}
}
