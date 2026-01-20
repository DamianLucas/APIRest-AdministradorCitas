package middleware

import (
	"adminApp/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequiereAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener header Authorization
		authHeader := c.GetHeader("Authorization")

		//verificar que exista
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token no proporcionado",
			})
			c.Abort()
			return
		}

		//verificar formato: "Bearer TOKEN"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		//validar token
		claims, err := utils.ValidarToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalido o expirado",
			})
			c.Abort()
			return
		}

		// Guardar claims en el contexto (disponible en handlers)
		c.Set("user_id", claims.UserID)
		c.Set("rol", claims.Rol)

		c.Next()
	}

}
