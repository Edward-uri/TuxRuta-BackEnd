package controller

import (
	"main/src/paradas/application"

	"github.com/gin-gonic/gin"
)

type DeleteParadaHandler struct {
	DeleteParadaUseCase *application.DeleteParadaUseCase
}

func NewDeleteParadaHandler(deleteParadaUseCase *application.DeleteParadaUseCase) *DeleteParadaHandler {
	return &DeleteParadaHandler{
		DeleteParadaUseCase: deleteParadaUseCase,
	}
}

func (c *DeleteParadaHandler) HandleDeleteParada(g *gin.Context) {
	var request struct {
		ID int `json:"id"`
	}

	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	err := c.DeleteParadaUseCase.Execute(request.ID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"message": "Parada deleted successfully"})
}
