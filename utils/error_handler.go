package utils

import "github.com/gin-gonic/gin"

func HandleServiceError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		Error(c, appErr.Status, appErr.Message)
		return
	}

	// fallback de seguridad
	InternalError(c)
}
