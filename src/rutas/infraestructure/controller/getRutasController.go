package controller

import (
	"main/src/rutas/application"

	"github.com/gin-gonic/gin"
)

type GetRutasHandler struct {
	GetRutasUseCase *application.GetRutasUseCase
}

func NewGetRutasHandler(getRutasUseCase *application.GetRutasUseCase) *GetRutasHandler {
	return &GetRutasHandler{
		GetRutasUseCase: getRutasUseCase,
	}
}

func (c *GetRutasHandler) HandleGetRutas(g *gin.Context) {
	rutas, err := c.GetRutasUseCase.Execute()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve rutas"})
		return
	}
	if len(rutas) == 0 {
		g.JSON(404, gin.H{"message": "No rutas found"})
		return
	}
	g.JSON(200, rutas)
}
