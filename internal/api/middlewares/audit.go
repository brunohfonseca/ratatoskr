package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuditMiddleware registra operações de modificação (POST, PUT, DELETE) no audit log
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Só audita operações de modificação
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Captura o body da request para salvar no audit
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// Restaura o body para o handler poder ler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Executa o handler
		c.Next()

		// Se a operação foi bem-sucedida (2xx), registra no audit
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go logAudit(c, bodyBytes) // Assíncrono para não impactar performance
		}
	}
}

// logAudit registra a operação no banco de dados
func logAudit(c *gin.Context, bodyBytes []byte) {
	// Extrai informações do contexto
	userID, exists := c.Get("id")
	if !exists {
		userID = nil // Operações sem autenticação
	}

	// Extrai a tabela/recurso do path
	tableName := extractTableName(c.FullPath())

	// Extrai o ID do recurso (se houver)
	recordID := c.Param("id")
	if recordID == "" {
		recordID = "unknown"
	}

	// Mapeia método HTTP para operação
	operation := mapHTTPMethodToOperation(c.Request.Method)

	// Prepara dados para salvar (body da requisição)
	var changedData map[string]interface{}
	if len(bodyBytes) > 0 {
		json.Unmarshal(bodyBytes, &changedData)
	} else {
		changedData = map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}
	}

	// Adiciona metadados extras
	changedData["_metadata"] = map[string]interface{}{
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"status":     c.Writer.Status(),
	}

	// Converte para JSONB
	changedDataJSON, err := json.Marshal(changedData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal audit data")
		return
	}

	// Insere no audit_log
	db := postgres.PostgresConn
	query := `
		INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(query, tableName, recordID, operation, changedDataJSON, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert audit log")
	}
}

// extractTableName extrai o nome da tabela/recurso do path
func extractTableName(path string) string {
	// Remove /api/v1/ e extrai o recurso principal
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			return parts[i+1] // Retorna a parte após /v1/
		}
	}

	// Fallback: retorna a primeira parte não vazia
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// mapHTTPMethodToOperation mapeia método HTTP para operação do audit
func mapHTTPMethodToOperation(method string) string {
	switch method {
	case "POST":
		return "INSERT"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	default:
		return method
	}
}
