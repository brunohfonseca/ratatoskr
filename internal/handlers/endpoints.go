package handlers

import (
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ListServices lista todos os endpoints cadastrados
func ListServices(c *gin.Context) {
	endpointService := services.NewEndpointService()

	endpoints, err := endpointService.List()
	if err != nil {
		log.Error().Err(err).Msg("Erro ao listar endpoints")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno do servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoints": endpoints,
		"total":     len(endpoints),
	})
}

// CreateService cria um novo endpoint
func CreateService(c *gin.Context) {
	var endpoint models.Endpoint

	// Bind JSON do request
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Criar endpoint
	endpointService := services.NewEndpointService()
	createdEndpoint, err := endpointService.CreateEndpoint(&endpoint)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar endpoint")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"endpoint": createdEndpoint,
		"message":  "Endpoint criado com sucesso",
	})
}

// GetService busca um endpoint específico por ID
func GetService(c *gin.Context) {
	id := c.Param("id")

	endpointService := services.NewEndpointService()
	endpoint, err := endpointService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoint": endpoint,
	})
}

// UpdateService atualiza um endpoint existente
func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var endpoint models.Endpoint

	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	endpointService := services.NewEndpointService()
	updatedEndpoint, err := endpointService.Update(id, &endpoint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpoint": updatedEndpoint,
		"message":  "Endpoint atualizado com sucesso",
	})
}

// DeleteService remove um endpoint
func DeleteService(c *gin.Context) {
	id := c.Param("id")

	endpointService := services.NewEndpointService()
	err := endpointService.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

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
