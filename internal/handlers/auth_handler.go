package handlers

import (
	"adminApp/internal/services"
	"adminApp/pkg/apperrors"
	"adminApp/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := services.Login(req.Email, req.Password)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp, "")
}
