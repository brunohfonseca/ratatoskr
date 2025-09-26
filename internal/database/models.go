package database

import "github.com/brunohfonseca/ratatoskr/internal/models"

// RegisterAllModels - Registra todas as models para sincronização automática
func RegisterAllModels() {
	RegisterModel(
		models.Endpoint{},
		models.EndPointHealthHistory{},
		models.AlertGroup{},
		models.AlertChannel{},
	)
}
