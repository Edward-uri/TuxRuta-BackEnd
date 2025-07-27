package controller

import (
	"main/src/paradas/application"
	"strconv"

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
	idParam := g.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		g.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	err = c.DeleteParadaUseCase.Execute(id)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"message": "Parada deleted successfully"})
}
