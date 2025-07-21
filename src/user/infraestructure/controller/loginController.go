package controller

import (
	"main/src/user/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	loginUseCase *application.LoginUseCase
}

func NewLoginHandler(loginUseCase *application.LoginUseCase) *LoginHandler {
	return &LoginHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *LoginHandler) HandleLogin(c *gin.Context) {
	var request application.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	response, err := h.loginUseCase.Execute(request)
	if err != nil {
		statusCode := http.StatusUnauthorized

		switch err.Error() {
		case "email is required", "password is required":
			statusCode = http.StatusBadRequest
		case "user account is disabled":
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}
