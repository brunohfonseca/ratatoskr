package monitors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
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
	checkSSLStr, _ := msg.Values["check_ssl"].(string)

	// Redis armazena valores como string, precisa converter para int
	expectedResponseCodeStr, _ := msg.Values["expected_response_code"].(string)
	expectedResponseCode, err := strconv.Atoi(expectedResponseCodeStr)
	if err != nil {
		logger.ErrLog("Erro ao converter expected_response_code", err)
		expectedResponseCode = 200 // Default
	}

	timeoutStr, _ := msg.Values["timeout"].(string)
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logger.ErrLog("Erro ao converter timeout", err)
		timeout = 30 // Default 30 segundos
	}

	url := fmt.Sprintf("%s%s", domain, path)
	check := doHealthCheck(url, expectedResponseCode, timeout)

	log := fmt.Sprintf("Checked Endpoint: UUID=%s, ExpectedResponseCode=%d, ResponseTime=%d, ResponseCode=%d, ResponseMessage=%s", uuid, expectedResponseCode, check.ResponseTime, check.ResponseStatusCode, check.ResponseMessage)
	logger.DebugLog(log)

	err = services.UpdateCheck(uuid, check)
	if err != nil {
		logger.ErrLog("Erro ao atualizar endpoint", err)
		return
	}

	// Adiciona o UUID ao resultado do check antes de registrar no histórico
	check.UUID = uuid

	err = services.RegisterCheck(check)
	if err != nil {
		logger.ErrLog("Erro ao registrar check no histórico", err)
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

func doHealthCheck(url string, expectedResponseCode, timeout int) models.EndpointResponse {
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
