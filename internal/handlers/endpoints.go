package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/brunohfonseca/ratatoskr/internal/entities"
	mongodb "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/mongodb"
)

// ListServices lista todos os endpoints cadastrados
func ListServices(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"endpoints": 0,
		"total":     0,
	})
}

// CreateService cria um novo endpoint
func CreateService(c *gin.Context) {
	// Ler o corpo uma única vez (permite múltiplos binds: map + entidade)
	var body map[string]interface{}
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Payload inválido",
			"details": err.Error(),
		})
		return
	}

	// Bind direto na sua entidade
	var ep entities.Endpoint
	if err := c.ShouldBindBodyWith(&ep, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Payload inválido",
			"details": err.Error(),
		})
		return
	}

	// Validações mínimas
	if ep.Name == "" || ep.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obrigatórios: name, domain"})
		return
	}

	// Compatibilidade: aceitar *_seconds (em segundos) e converter para time.Duration
	if v, ok := body["timeout_seconds"].(float64); ok && v > 0 {
		ep.Timeout = time.Duration(int64(v)) * time.Second
	}
	if v, ok := body["interval_seconds"].(float64); ok && v > 0 {
		ep.Interval = time.Duration(int64(v)) * time.Second
	}

	// Defaults controlados pelo servidor
	if ep.Timeout == 0 {
		ep.Timeout = 30 * time.Second
	}
	if ep.Interval == 0 {
		ep.Interval = 5 * time.Minute
	}
	ep.Status = entities.StatusUnknown
	ep.CreatedAt = time.Now().UTC()
	ep.UpdatedAt = time.Now().UTC()

	// Inserção no Mongo usando a conexão global
	if mongodb.MongoDatabase == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB não inicializado"})
		return
	}
	coll := mongodb.MongoDatabase.Collection("endpoints")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := coll.InsertOne(ctx, ep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao criar endpoint",
			"details": err.Error(),
		})
		return
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		ep.ID = oid
	}

	c.JSON(http.StatusCreated, gin.H{
		"endpoint": ep,
		"message":  "Endpoint criado com sucesso",
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
