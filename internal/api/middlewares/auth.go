package middlewares

import (
	"net/http"
	"strings"

	"github.com/brunohfonseca/ratatoskr/internal/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida o token JWT nas requisições
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação não fornecido"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		// Adicionar informações do usuário no contexto
		c.Set("uuid", claims.UserUUID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
