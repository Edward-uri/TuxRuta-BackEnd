package controllers

import (
	"main/src/user/application"

	"github.com/gin-gonic/gin"
)

type GetUsersHandler struct {
	GetUsersUseCase *application.GetAllUserUseCase
}

func NewGetUsersHandler(getUsersUseCase *application.GetAllUserUseCase) *GetUsersHandler {
	return &GetUsersHandler{
		GetUsersUseCase: getUsersUseCase,
	}
}

func (c *GetUsersHandler) HandleGetUsers(g *gin.Context) {
	users, err := c.GetUsersUseCase.Execute()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve users"})
		return
	}
	if len(users) == 0 {
		g.JSON(404, gin.H{"message": "No users found"})
	}
	g.JSON(200, users)
}
