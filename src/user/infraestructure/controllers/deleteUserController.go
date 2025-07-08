package controllers

import (
	"main/src/user/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteUserHandler struct {
	DeleteUserUseCase *application.DeleteUserUseCase
}

func NewDeleteUserHandler(deleteUserHandler *application.DeleteUserUseCase) *DeleteUserHandler {
	return &DeleteUserHandler{
		DeleteUserUseCase: deleteUserHandler,
	}
}

func (c *DeleteUserHandler) HandleDeleteUser(g *gin.Context) {
	userIdParam := g.Param("id")
	id, err := strconv.Atoi(userIdParam)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	err = c.DeleteUserUseCase.Execute(id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
