package services

import (
	"context"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	infraRedis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/redis/go-redis/v9"
)

// CreateEndpoint salva o endpoint no banco e envia pro Redis
func CreateEndpoint(endpoint *models.Endpoint, userID int) error {
	db := postgres.PostgresConn

	sql := "INSERT INTO endpoints (name, domain, path, check_ssl, last_modified_by) VALUES ($1, $2, $3, $4, $5) RETURNING uuid, status"
	err := db.QueryRow(sql,
		endpoint.Name,
		endpoint.Domain,
		endpoint.EndpointPath,
		endpoint.CheckSSL,
		userID,
	).Scan(&endpoint.UUID, &endpoint.Status)

	if err != nil {
		return err
	}

	// Envia pro Redis Stream
	ctx := context.Background()
	redisClient := infraRedis.RedisClient

	_, err = redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "endpoints",
		Values: map[string]interface{}{
			"uuid":      endpoint.UUID,
			"domain":    endpoint.Domain,
			"path":      endpoint.EndpointPath,
			"check_ssl": endpoint.CheckSSL,
		},
	}).Result()

	return err
}

func UpdateCheck(endpoint models.Endpoint) {
	//db := postgres.PostgresConn

	//sql := "UPDATE endpoints SET status = $1 WHERE uuid = $2"
}

func GetEndpointByUUID(uuid string) (interface{}, error) {
	db := postgres.PostgresConn

	sql := `
		SELECT 
			uuid,
			expected_response_code,
			timeout_seconds, 
			alert_group_id
		FROM endpoints
		WHERE uuid = $1
	`
	row := db.QueryRow(sql, uuid)
	var endpoint models.Endpoint
	err := row.Scan(&endpoint.UUID, &endpoint.ExpectedResponseCode, &endpoint.TimeoutSeconds, &endpoint.AlertGroupID)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}
