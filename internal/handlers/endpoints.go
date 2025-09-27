package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/brunohfonseca/ratatoskr/internal/config"
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
	// Payload de criação (simples, com defaults de timeout/intervalo)
	type createEndpointRequest struct {
		Name            string      `json:"name" binding:"required"`
		Domain          string      `json:"domain" binding:"required"`
		Port            int         `json:"port"`
		Endpoint        string      `json:"endpoint"`
		TimeoutSeconds  int         `json:"timeout_seconds"`
		IntervalSeconds int         `json:"interval_seconds"`
		CheckSSL        bool        `json:"check_ssl"`
		Authentication  interface{} `json:"authentication"`
		AlertGroupIDs   []string    `json:"alert_group_ids"`
		Enabled         *bool       `json:"enabled"`
	}

	var req createEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Payload inválido",
			"details": err.Error(),
		})
		return
	}

	// Defaults
	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	interval := time.Duration(req.IntervalSeconds) * time.Second
	if interval == 0 {
		interval = 5 * time.Minute
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	// Converter AlertGroupIDs
	var alertIDs []primitive.ObjectID
	for _, s := range req.AlertGroupIDs {
		id, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "alert_group_ids contém ID inválido",
				"details": s,
			})
			return
		}
		alertIDs = append(alertIDs, id)
	}

	// Montar entidade para persistência
	ep := entities.Endpoint{
		Name:          req.Name,
		Domain:        req.Domain,
		Port:          req.Port,
		Endpoint:      req.Endpoint,
		Timeout:       timeout,
		Interval:      interval,
		CheckSSL:      req.CheckSSL,
		Status:        entities.StatusUnknown,
		Enabled:       enabled,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		AlertGroupIDs: alertIDs,
	}
	ep.Authentication = req.Authentication

	// Obter DB/coleção a partir da URL do Mongo configurada
	cfg := config.Get()
	if cfg == nil || cfg.Database.MongoURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuração do MongoDB ausente"})
		return
	}
	cs, err := connstring.ParseAndValidate(cfg.Database.MongoURL)
	if err != nil || cs.Database == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuração do MongoDB inválida: defina o database na URL"})
		return
	}
	coll := mongodb.MongoClient.Database(cs.Database).Collection("endpoints")

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
