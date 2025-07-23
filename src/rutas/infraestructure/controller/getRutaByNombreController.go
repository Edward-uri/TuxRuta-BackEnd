package controller

import (
	"main/src/rutas/application"

	"github.com/gin-gonic/gin"
)

type GetRutaByNombreHandler struct {
	GetRutaByNombreUseCase *application.GetRutaByNombreUseCase
}

func NewGetRutaByNombreHandler(getRutaByNombreUseCase *application.GetRutaByNombreUseCase) *GetRutaByNombreHandler {
	return &GetRutaByNombreHandler{
		GetRutaByNombreUseCase: getRutaByNombreUseCase,
	}
}

func (c *GetRutaByNombreHandler) HandleGetRutaByNombre(g *gin.Context) {
	nombre := g.Param("nombre")
	if nombre == "" {
		g.JSON(400, gin.H{"error": "Nombre is required"})
		return
	}
	ruta, err := c.GetRutaByNombreUseCase.Execute(nombre)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve ruta"})
		return
	}
	g.JSON(200, ruta)
}
