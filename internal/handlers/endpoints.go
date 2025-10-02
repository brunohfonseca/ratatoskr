package handlers

import (
	"database/sql"
	"errors"
	"net/http"

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
		logger.DebugLog(err.Error())
		msg := errors.New("possible missing required fields, check your payload")
		responses.Error(c, http.StatusBadRequest, msg)
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

	response, err := services.ListEndpoints()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}

	responses.Success(c, http.StatusOK, gin.H{
		"total":     len(response),
		"endpoints": response,
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

// CheckEndpoint adiciona um endpoint a fila de verificação
func CheckEndpoint(c *gin.Context) {
	var endpoint struct {
		UUID string `json:"uuid,omitempty" binding:"required"`
	}

	if err := c.BindJSON(&endpoint); err != nil {
		logger.DebugLog(err.Error())
		msg := errors.New("missing uuid")
		responses.Error(c, http.StatusBadRequest, msg)
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
			"uuid":    ep.UUID,
			"name":    ep.Name,
			"domain":  ep.Domain,
			"path":    ep.EndpointPath,
			"timeout": ep.Timeout,
		},
	})
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}
	if ep.CheckSSL {
		err := infraRedis.StreamPublish(c, &redis.XAddArgs{
			Stream: "ssl-checks",
			Values: map[string]interface{}{
				"uuid":    ep.UUID,
				"domain":  ep.Domain,
				"timeout": ep.Timeout,
			},
		})
		logger.DebugLog("SSL check adicionado a fila de verificação")
		if err != nil {
			responses.Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	responses.Success(c, http.StatusOK, gin.H{
		"uuid":    ep.UUID,
		"name":    ep.Name,
		"message": "Endpoint adicionado a fila de verificação",
	})
}
