package handlers

import (
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HealthCheck verifica se a aplicação está funcionando
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Ratatoskr API está funcionando",
		"service": "ratatoskr-api",
	})
}

// ReadinessCheck verifica se a aplicação está pronta para receber tráfego
func ReadinessCheck(c *gin.Context) {
	// Verificar conexões com dependências
	mongoHealthy, mongoStatus, mongoErr := database.CheckMongoDBHealth()
	redisHealthy, redisStatus, redisErr := database.CheckRedisHealth()

	// Determinar status geral da aplicação
	overallHealthy := mongoHealthy && redisHealthy
	httpStatus := http.StatusOK

	if !overallHealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	// Preparar resposta detalhada
	checks := gin.H{
		"mongodb": gin.H{
			"status":  mongoStatus,
			"healthy": mongoHealthy,
		},
		"redis": gin.H{
			"status":  redisStatus,
			"healthy": redisHealthy,
		},
	}

	// Adicionar detalhes de erro se houver
	if mongoErr != nil {
		checks["mongodb"].(gin.H)["error"] = mongoErr.Error()
		log.Warn().Err(mongoErr).Msg("MongoDB health check failed")
	}

	if redisErr != nil {
		checks["redis"].(gin.H)["error"] = redisErr.Error()
		log.Warn().Err(redisErr).Msg("Redis health check failed")
	}

	status := "ready"
	if !overallHealthy {
		status = "not_ready"
	}

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"checks":  checks,
		"healthy": overallHealthy,
	})
}

// LivenessCheck verifica se a aplicação ainda está viva
func LivenessCheck(c *gin.Context) {
	log.Debug().Msg("Liveness check executado")

	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": gin.H{},
	})
}
