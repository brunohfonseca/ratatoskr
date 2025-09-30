package handlers

import (
	"net/http"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateService cria um novo endpoint
func CreateService(c *gin.Context) {
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

	v7, err := uuid.NewV7()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}

	db := postgres.PostgresConn
	sql := "INSERT INTO endpoints (name, uuid, domain, path, check_ssl, last_modified_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = db.QueryRow(sql,
		endpoint.Name,
		v7,
		endpoint.Domain,
		endpoint.EndpointPath,
		endpoint.CheckSSL,
		userID,
	).Scan(&endpoint.ID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}
	responses.Success(c, http.StatusCreated, endpoint)
}

// ListServices lista todos os endpoints cadastrados
func ListServices(c *gin.Context) {
	var endpoints []models.Endpoint

	db := postgres.PostgresConn
	sql := `
		SELECT
		    id,
		    uuid,
		    name,
		    domain,
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

// GetService busca um endpoint específico por ID
func GetService(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"endpoint": id,
	})
}

// UpdateService atualiza um endpoint existente
func UpdateService(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"endpoint": id,
		"message":  "Endpoint atualizado com sucesso",
	})
}

// DeleteService remove um endpoint
func DeleteService(c *gin.Context) {
	//id := c.Param("id")

	// TODO: Implementar health check real
	c.JSON(http.StatusOK, gin.H{
		"message": "Endpoint removido com sucesso",
	})
}

// GetServiceStatus retorna o status atual do serviço
func GetServiceStatus(c *gin.Context) {
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

// GetServiceHistory retorna o histórico de health checks
func GetServiceHistory(c *gin.Context) {
	id := c.Param("id")

	// TODO: Buscar histórico no MongoDB
	c.JSON(http.StatusOK, gin.H{
		"service_id": id,
		"history":    []gin.H{},
		"total":      0,
		"message":    "Histórico do serviço (implementação pendente)",
	})
}

// GetServiceUptime retorna estatísticas de uptime
func GetServiceUptime(c *gin.Context) {
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
