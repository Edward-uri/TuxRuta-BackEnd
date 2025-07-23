package controller

import (
	"main/src/paradas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetParadaHandler struct {
	GetParadaUseCase *application.GetParadaUseCase
}

func NewGetParadaHandler(getParadaUseCase *application.GetParadaUseCase) *GetParadaHandler {
	return &GetParadaHandler{
		GetParadaUseCase: getParadaUseCase,
	}
}

func (c *GetParadaHandler) HandleGetParada(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid parada ID"})
		return
	}

	parada, err := c.GetParadaUseCase.Execute(id)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"parada": parada})
}
