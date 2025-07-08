package controllers

import (
	"main/src/user/application"
	"main/src/user/domain/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUserByIDHandler struct {
	GetUserByIDUseCase *application.GetUserByIdUseCase
}

func NewGetUserByIDHandler(getUserByIDUseCase *application.GetUserByIdUseCase) *GetUserByIDHandler {
	return &GetUserByIDHandler{
		GetUserByIDUseCase: getUserByIDUseCase,
	}
}

func (c *GetUserByIDHandler) HandleGetUserByID(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.GetUserByIDUseCase.Execute(id)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if (user == entities.Usuario{}) {
		g.JSON(404, gin.H{"message": "User not found"})
		return
	}

	g.JSON(200, user)
}
