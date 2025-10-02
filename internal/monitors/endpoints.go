package monitors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/services"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/redis/go-redis/v9"
)

func ProcessEndpoint(ctx context.Context, redisClient *redis.Client, stream, group string, msg redis.XMessage) {
	// Leitura segura dos campos do Redis
	uuid, _ := msg.Values["uuid"].(string)
	domain, _ := msg.Values["domain"].(string)
	path, _ := msg.Values["path"].(string)
	expectedResponseCode, _ := msg.Values["expected_response_code"].(int)
	timeout, _ := msg.Values["timeout"].(int)
	checkSSLStr, _ := msg.Values["check_ssl"].(string)

	url := fmt.Sprintf("%s%s", domain, path)
	check := doHealthCheck(url, expectedResponseCode, timeout)

	log := fmt.Sprintf("Checked Endpoint: UUID=%s, ExpectedResponseCode=%d, ResponseTime=%d, ResponseCode=%d, ResponseMessage=%s", check.UUID, check.ExpectedResponseCode, check.ResponseTime, check.ResponseStatusCode, check.ResponseMessage)
	logger.DebugLog(log)

	err := services.UpdateCheck(uuid, check)
	if err != nil {
		logger.ErrLog("Erro ao atualizar endpoint", err)
		return
	}

	if checkSSLStr == "true" {
		_, err := FetchSSL(domain)
		if err != nil {
			return
		}
	}

	logMsg := fmt.Sprintf("✅ Health check completed in %s", uuid)
	logger.DebugLog(logMsg)
}

func doHealthCheck(url string, timeout, expectedResponseCode int) models.EndpointResponse {
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
			Status:             models.StatusOffline,
			ResponseTime:       int(time.Since(start).Milliseconds()),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedResponseCode {
		return models.EndpointResponse{
			ResponseStatusCode: resp.StatusCode,
			ResponseMessage:    "Response code does not match expected code",
			Status:             models.StatusError,
			ResponseTime:       int(time.Since(start).Milliseconds()),
		}
	}

	return models.EndpointResponse{
		ResponseStatusCode: resp.StatusCode,
		ResponseMessage:    "Success",
		Status:             models.StatusOnline,
		ResponseTime:       int(time.Since(start).Milliseconds()),
	}
}
