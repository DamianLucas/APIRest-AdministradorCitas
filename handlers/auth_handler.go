package handlers

import (
	"adminApp/services"
	"adminApp/utils"
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
		utils.BadRequest(c, "datos inválidos")
		return
	}

	resp, err := services.Login(req.Email, req.Password)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, resp, "")
}
