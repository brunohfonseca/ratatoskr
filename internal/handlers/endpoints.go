package handlers

import (
	"net/http"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/gin-gonic/gin"
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

// GetEndpointStatus retorna o status atual do serviço
func GetEndpointStatus(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implementar verificação real de status
	c.JSON(http.StatusOK, gin.H{
		"service_id":    id,
		"status":        "online",
		"last_check":    "2024-01-01T00:00:00Z",
		"response_time": 120,
		"message":       "Status do serviço (implementação pendente)",
	})
}

// TriggerHealthCheck força uma verificação de health check
func TriggerHealthCheck(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implementar health check real
	c.JSON(http.StatusOK, gin.H{
		"service_id": id,
		"status":     "check_triggered",
		"message":    "Health check iniciado (implementação pendente)",
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
