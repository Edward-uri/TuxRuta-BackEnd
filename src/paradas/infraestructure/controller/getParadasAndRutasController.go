package controller

import (
	"main/src/paradas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetParadasAndRutasHandler struct {
	GetParadasAndRutasUseCase *application.GetParadasAndRutaUseCase
}

func NewGetParadasAndRutasHandler(getParadasAndRutasUseCase *application.GetParadasAndRutaUseCase) *GetParadasAndRutasHandler {
	return &GetParadasAndRutasHandler{
		GetParadasAndRutasUseCase: getParadasAndRutasUseCase,
	}
}

func (c *GetParadasAndRutasHandler) HandleGetParadasAndRutas(g *gin.Context) {
	rutaIDStr := g.Param("ruta_id")
	rutaID, err := strconv.Atoi(rutaIDStr)
	if err != nil || rutaID <= 0 {
		g.JSON(400, gin.H{"error": "Invalid ruta_id"})
		return
	}

	paradas, err := c.GetParadasAndRutasUseCase.Execute(rutaID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"paradas": paradas})
}
