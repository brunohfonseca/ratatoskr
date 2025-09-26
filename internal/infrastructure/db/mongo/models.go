package mongo

import (
	"github.com/brunohfonseca/ratatoskr/internal/models"
)

// RegisterAllModels - Registra todas as models para sincronização automática
func RegisterAllModels() {
	RegisterModel(
		models.Endpoint{},
		models.EndpointHealthHistory{},
		models.AlertGroup{},
		models.AlertChannel{},
	)
}
