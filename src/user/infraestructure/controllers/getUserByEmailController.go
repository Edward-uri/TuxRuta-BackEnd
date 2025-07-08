package controllers

import (
	"main/src/user/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetUserByEmailHandler struct {
	GetUserByEmailUseCase *application.GetUserByEmailUseCase
}

func NewGetUserByEmailHandler(getUserByEmailUseCase *application.GetUserByEmailUseCase) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{
		GetUserByEmailUseCase: getUserByEmailUseCase,
	}
}

func (c *GetUserByEmailHandler) HandleGetUserByEmail(g *gin.Context) {
	email, exists := g.GetQuery("email")
	if !exists || email == "" {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": "Email parameter is required",
		})
		return
	}

	if !isValidEmail(email) {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	user, err := c.GetUserByEmailUseCase.Execute(email)
	if err != nil {
		if err.Error() == "user not found" {
			g.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		g.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	if user.ID == 0 {
		g.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user.Password = ""

	g.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func isValidEmail(email string) bool {
	return len(email) > 0 &&
		len(email) <= 254 &&
		contains(email, "@") &&
		contains(email, ".")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
