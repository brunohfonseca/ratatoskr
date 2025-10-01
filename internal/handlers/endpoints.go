package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	infraRedis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/services"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CreateEndpoint cria um novo endpoint
func CreateEndpoint(c *gin.Context) {
	var endpoint models.Endpoint
	if err := c.BindJSON(&endpoint); err != nil {
		responses.Error(c, http.StatusBadRequest, err)
		return
	}

	// Extrair user_id do contexto (colocado pelo middleware JWT)
	userIDInterface, exists := c.Get("id")
	if !exists {
		responses.ErrorMsg(c, http.StatusInternalServerError, "Usuário não autenticado")
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		responses.ErrorMsg(c, http.StatusInternalServerError, "ID de usuário inválido")
		return
	}

	// Chama o service
	if err := services.CreateEndpoint(&endpoint, userID); err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}

	responses.Success(c, http.StatusCreated, endpoint)
}

// ListEndpoints lista todos os endpoints cadastrados
func ListEndpoints(c *gin.Context) {
	var endpoints []models.Endpoint

	db := postgres.PostgresConn
	sql := `
		SELECT
		    id,
		    uuid,
		    name,
		    domain,
		    status,
		    path,
		    check_ssl,
		    last_modified_by 
		FROM endpoints`
	rows, err := db.Query(sql)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	// Iterar sobre as rows e popular o slice
	for rows.Next() {
		var endpoint models.Endpoint
		var lastModifiedBy *int

		err := rows.Scan(
			&endpoint.ID,
			&endpoint.UUID,
			&endpoint.Name,
			&endpoint.Domain,
			&endpoint.Status,
			&endpoint.EndpointPath,
			&endpoint.CheckSSL,
			&lastModifiedBy,
		)
		if err != nil {
			responses.Error(c, http.StatusInternalServerError, err)
			return
		}

		endpoints = append(endpoints, endpoint)
	}

	// Verificar se houve erro durante a iteração
	if err = rows.Err(); err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}

	responses.Success(c, http.StatusOK, gin.H{
		"total":     len(endpoints),
		"endpoints": endpoints,
	})
}

// GetEndpoint busca um endpoint específico por ID
func GetEndpoint(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"endpoint": id,
	})
}

// UpdateEndpoint atualiza um endpoint existente
func UpdateEndpoint(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"endpoint": id,
		"message":  "Endpoint atualizado com sucesso",
	})
}

// DeleteEndpoint remove um endpoint
func DeleteEndpoint(c *gin.Context) {
	//id := c.Param("id")

	// TODO: Implementar health check real
	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint removido com sucesso",
	})
}

// GetEndpointHistory retorna o histórico de health checks
func GetEndpointHistory(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"service_id": id,
		"history":    []gin.H{},
		"total":      0,
		"message":    "Histórico do serviço (implementação pendente)",
	})
}

// GetEndpointUptime retorna estatísticas de uptime
func GetEndpointUptime(c *gin.Context) {
	id := c.Param("id")

	// TODO: Calcular uptime real baseado no histórico
	c.JSON(http.StatusOK, gin.H{
		"service_id": id,
		"uptime": gin.H{
			"percentage":        99.9,
			"total_checks":      1000,
			"successful_checks": 999,
			"failed_checks":     1,
		},
		"message": "Estatísticas de uptime (implementação pendente)",
	})
}

func CheckEndpoint(c *gin.Context) {
	var endpoint models.Endpoint
	if err := c.BindJSON(&endpoint); err != nil {
		responses.Error(c, http.StatusBadRequest, err)
		return
	}

	if endpoint.UUID == "" {
		responses.ErrorMsg(c, http.StatusBadRequest, "UUID do endpoint não informado")
		return
	}

	ep, err := services.GetEndpointByUUID(endpoint.UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responses.ErrorMsg(c, http.StatusNotFound, "Endpoint não encontrado")
		} else {
			responses.Error(c, http.StatusInternalServerError, err)
		}
		return
	}

	logger.DebugLog("Endpoint localizado: " + ep.Name)

	err = infraRedis.StreamPublish(c, &redis.XAddArgs{
		Stream: "endpoints",
		Values: map[string]interface{}{
			"uuid":      endpoint.UUID,
			"name":      endpoint.Name,
			"domain":    endpoint.Domain,
			"path":      endpoint.EndpointPath,
			"timeout":   endpoint.Timeout,
			"check_ssl": endpoint.CheckSSL,
		},
	})
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}

	responses.Success(c, http.StatusOK, gin.H{
		"uuid":    ep.UUID,
		"message": "Endpoint adicionado a fila de verificação",
	})
}
