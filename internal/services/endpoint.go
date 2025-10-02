package services

import (
	"context"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	infraRedis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// CreateEndpoint salva o endpoint no banco e envia pro Redis
func CreateEndpoint(endpoint *models.Endpoint, userID string) error {
	db := postgres.PostgresConn

	sql := "INSERT INTO endpoints (name, domain, path, check_ssl, expected_response_code, timeout_seconds, last_modified_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING uuid, status"
	err := db.QueryRow(sql,
		endpoint.Name,
		endpoint.Domain,
		endpoint.EndpointPath,
		endpoint.CheckSSL,
		endpoint.ExpectedResponseCode,
		endpoint.TimeoutSeconds,
		userID,
	).Scan(&endpoint.UUID, &endpoint.Status)
	if err != nil {
		return err
	}

	// Envia pro Redis Stream
	ctx := context.Background()
	err = infraRedis.StreamPublish(ctx, &redis.XAddArgs{
		Stream: "endpoints",
		Values: map[string]interface{}{
			"uuid":                   endpoint.UUID,
			"name":                   endpoint.Name,
			"domain":                 endpoint.Domain,
			"path":                   endpoint.EndpointPath,
			"timeout":                endpoint.Timeout,
			"expected_response_code": endpoint.ExpectedResponseCode,
			"check_ssl":              endpoint.CheckSSL,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("⚠️ Erro ao publicar no Redis")
		return err
	}

	return err
}

func UpdateCheck(uuid string, endpoint models.EndpointResponse) error {
	db := postgres.PostgresConn

	sql := `
		UPDATE 
		    endpoints 
		SET 
		    status = $2,
		    response_time_ms = $3,
		    response_code = $4,
		    response_message = $5
		WHERE uuid = $1
		`
	_, err := db.Exec(sql,
		uuid,
		endpoint.Status,
		endpoint.ResponseTime,
		endpoint.ResponseStatusCode,
		endpoint.ResponseMessage,
	)
	if err != nil {
		return err
	}
	return nil
}

func RegisterCheck(endpoint models.EndpointResponse) error {
	db := postgres.PostgresConn

	query := `
		INSERT INTO
			endpoint_checks(endpoint_id, status, response_time_ms, response_message)
		VALUES
			($1, $2, $3, $4)
		`
	_, err := db.Exec(query,
		endpoint.UUID,
		endpoint.Status,
		endpoint.ResponseTime,
		endpoint.ResponseMessage,
	)
	if err != nil {
		logger.ErrLog("Erro ao registrar check", err)
		return err
	}
	return nil
}

func GetEndpointByUUID(uuid string) (models.Endpoint, error) {
	var endpoint models.Endpoint
	db := postgres.PostgresConn

	sql := `
		SELECT 
			uuid,
			name,
			expected_response_code,
			check_ssl,
			timeout_seconds, 
			alert_group_id
		FROM endpoints
		WHERE uuid = $1
	`
	err := db.QueryRow(sql, uuid).Scan(
		&endpoint.UUID,
		&endpoint.Name,
		&endpoint.ExpectedResponseCode,
		&endpoint.CheckSSL,
		&endpoint.TimeoutSeconds,
		&endpoint.AlertGroupID,
	)
	if err != nil {
		return models.Endpoint{}, err
	}

	return endpoint, nil
}

func ListEndpoints() ([]models.Endpoint, error) {
	var endpoints []models.Endpoint
	db := postgres.PostgresConn

	query := `
		SELECT
		    uuid,
		    name,
		    domain,
		    status,
		    check_ssl
		FROM endpoints
		`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterar sobre as rows e popular o slice
	for rows.Next() {
		var endpoint models.Endpoint

		err := rows.Scan(
			&endpoint.UUID,
			&endpoint.Name,
			&endpoint.Domain,
			&endpoint.Status,
			&endpoint.CheckSSL,
		)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}
