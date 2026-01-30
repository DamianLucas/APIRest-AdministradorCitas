package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequiereRol(rolesPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener rol desde el contexto (set en el JWT middleware)
		rol, exists := c.Get("rol")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Rol no encontrado",
			})
			c.Abort()
			return
		}

		rolStr, ok := rol.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Rol invalido",
			})
			c.Abort()
			return
		}

		// Verificar si el rol está permitido
		for _, permitido := range rolesPermitidos {
			if rolStr == permitido {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "No tienes permisos para esta accion",
		})
		c.Abort()
	}
}
