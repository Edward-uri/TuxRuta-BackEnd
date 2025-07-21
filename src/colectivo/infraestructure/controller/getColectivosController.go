package controller

import (
	"main/src/colectivo/application"

	"github.com/gin-gonic/gin"
)

type GetColectivosHandler struct {
	GetColectivosUseCase *application.GetColectivosUseCase
}

func NewGetColectivosHandler(getColectivosUseCase *application.GetColectivosUseCase) *GetColectivosHandler {
	return &GetColectivosHandler{
		GetColectivosUseCase: getColectivosUseCase,
	}
}

func (c *GetColectivosHandler) HandleGetColectivos(g *gin.Context) {
	colectivos, err := c.GetColectivosUseCase.Execute()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve colectivos"})
		return
	}
	if len(colectivos) == 0 {
		g.JSON(404, gin.H{"message": "No colectivos found"})
		return
	}
	g.JSON(200, colectivos)
}
