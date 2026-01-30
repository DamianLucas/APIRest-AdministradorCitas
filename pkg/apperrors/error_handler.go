package apperrors

import (
	"adminApp/pkg/response"

	"github.com/gin-gonic/gin"
)

func HandleServiceError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		response.Error(c, appErr.Status, appErr.Message)
		return
	}

	// fallback de seguridad
	response.InternalError(c)
}
