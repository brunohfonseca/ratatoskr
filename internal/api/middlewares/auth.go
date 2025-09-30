package middlewares

import (
	"net/http"
	"strings"

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

		// Adicionar informações do usuário no contexto
		c.Set("id", claims.UserID)
		c.Set("uuid", claims.UserUUID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
