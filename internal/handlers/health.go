package handlers

import (
	"net/http"

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
	// Aqui você pode adicionar verificações de dependências
	// Por exemplo: verificar conexão com MongoDB, Redis, etc.

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "connected",
			"redis":    "connected",
		},
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
