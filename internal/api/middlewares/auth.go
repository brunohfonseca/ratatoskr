package middlewares

import (
	"net/http"
	"strings"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/utils"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida o token JWT nas requisições
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responses.ErrorMsg(c, http.StatusUnauthorized, "Missing Authorization Bearer token")
			c.Abort()
			return
		}

		// Remover prefixo "Bearer " se existir
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Se não tinha "Bearer ", tentar sem prefixo
			tokenString = authHeader
		}

		// Validar token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			responses.ErrorMsg(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Validar se usuário ainda existe e dados batem com o banco
		db := postgres.PostgresConn
		var userID int
		var userUUID, userEmail string
		var enabled bool

		err = db.QueryRow(
			"SELECT id, uuid, email, enabled FROM users WHERE id = $1",
			claims.UserID,
		).Scan(&userID, &userUUID, &userEmail, &enabled)

		if err != nil {
			responses.ErrorMsg(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}

		// Verifica se UUID do token bate com o do banco
		if userUUID != claims.UserUUID {
			responses.ErrorMsg(c, http.StatusUnauthorized, "Token invalidated - user data changed")
			c.Abort()
			return
		}

		// Verifica se email do token bate com o do banco
		if userEmail != claims.Email {
			responses.ErrorMsg(c, http.StatusUnauthorized, "Token invalidated - user data changed")
			c.Abort()
			return
		}

		// Verifica se usuário está habilitado
		if !enabled {
			responses.ErrorMsg(c, http.StatusUnauthorized, "User is disabled")
			c.Abort()
			return
		}

		// Adicionar informações do usuário no contexto
		c.Set("id", claims.UserID)
		c.Set("uuid", claims.UserUUID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
