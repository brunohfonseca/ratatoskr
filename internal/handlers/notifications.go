package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ListAlerts lista todas as configurações de alerta
func ListAlerts(c *gin.Context) {
	log.Info().Msg("Listando configurações de alerta")

	c.JSON(http.StatusOK, gin.H{
		"alerts":  []gin.H{},
		"total":   0,
		"message": "Lista de alertas (implementação pendente)",
	})
}

// CreateAlert cria uma nova configuração de alerta
func CreateAlert(c *gin.Context) {
	log.Info().Msg("Criando nova configuração de alerta")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Alerta criado com sucesso (implementação pendente)",
		"id":      "temp-alert-id",
	})
}

// UpdateAlert atualiza uma configuração de alerta existente
func UpdateAlert(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("alert_id", id).Msg("Atualizando configuração de alerta")

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerta atualizado com sucesso (implementação pendente)",
		"id":      id,
	})
}

// DeleteAlert remove uma configuração de alerta
func DeleteAlert(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("alert_id", id).Msg("Removendo configuração de alerta")

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerta removido com sucesso (implementação pendente)",
		"id":      id,
	})
}

// GetAlertsHistory retorna o histórico de alertas disparados
func GetAlertsHistory(c *gin.Context) {
	log.Info().Msg("Buscando histórico de alertas")

	// Parâmetros opcionais
	limit := c.DefaultQuery("limit", "50")

	c.JSON(http.StatusOK, gin.H{
		"alerts":  []gin.H{},
		"total":   0,
		"limit":   limit,
		"message": "Histórico de alertas (implementação pendente)",
	})
}

// TestAlert testa uma configuração de alerta específica
func TestAlert(c *gin.Context) {
	id := c.Param("id")
	log.Info().Str("alert_id", id).Msg("Testando configuração de alerta")

	// TODO: Implementar teste real de alerta
	c.JSON(http.StatusOK, gin.H{
		"alert_id": id,
		"status":   "test_sent",
		"message":  "Teste de alerta enviado (implementação pendente)",
	})
}
