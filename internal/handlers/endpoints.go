package handlers

import (
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/repositories"
	"github.com/gin-gonic/gin"
)

type EndpointHandler struct {
	repo repositories.EndpointRepository
}

func NewEndpointHandler(repo repositories.EndpointRepository) *EndpointHandler {
	return &EndpointHandler{repo: repo}
}

// CreateService cria um novo endpoint
func CreateService(c *gin.Context) {
	//var e entities.Endpoint
	//if err := c.ShouldBindJSON(&e); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
	//	return
	//}
	//
	//// valida campos obrigatórios
	//if e.Name == "" || e.Domain == "" {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Name e Domain são obrigatórios"})
	//	return
	//}
	//
	//// contexto com timeout baseado no request
	//ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	//defer cancel()
	//
	//id, err := h.repo.Create(ctx, &e)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//e.ID = id // garante que o JSON de resposta tem o ID gerado
	//c.JSON(http.StatusCreated, gin.H{
	//	"endpoint": e,
	//})
}

// ListServices lista todos os endpoints cadastrados
func ListServices(c *gin.Context) {
	//ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	//defer cancel()
	//
	//endpoints, err := h.repo.FindAll(ctx)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"total":     len(endpoints),
	//	"endpoints": endpoints,
	//})
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
