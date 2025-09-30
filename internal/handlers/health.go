package handlers

import (
	"net/http"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
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
	postgresHealthy, postgresStatus, postgresErr := postgres.CheckPostgresHealth()
	redisHealthy, redisStatus, redisErr := redis.CheckRedisHealth()

	// Determinar status geral da aplicação
	overallHealthy := postgresHealthy && redisHealthy
	httpStatus := http.StatusOK

	if !overallHealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	// Preparar resposta detalhada
	checks := gin.H{
		"postgres": gin.H{
			"status":  postgresStatus,
			"healthy": postgresHealthy,
		},
		"redis": gin.H{
			"status":  redisStatus,
			"healthy": redisHealthy,
		},
	}

	// Adicionar detalhes de erro se houver
	if postgresErr != nil {
		checks["postgres"].(gin.H)["error"] = postgresErr.Error()
		log.Warn().Err(postgresErr).Msg("Postgres health check failed")
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
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": gin.H{},
	})
}
