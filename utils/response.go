package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, data interface{}, message string) {
	resp := gin.H{
		"success": true,
	}

	if data != nil {
		resp["data"] = data
	}

	if message != "" {
		resp["message"] = message
	}
	c.JSON(status, resp)
}

func Error(c *gin.Context, status int, err string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   err,
	})
}

func BadRequest(c *gin.Context, err string) {
	Error(c, http.StatusBadRequest, err)
}

func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "Error interno del servidor")
}

func Forbidden(c *gin.Context, err string) {
	Error(c, http.StatusForbidden, err)
}

func NotFound(c *gin.Context, err string) {
	Error(c, http.StatusNotFound, err)
}
