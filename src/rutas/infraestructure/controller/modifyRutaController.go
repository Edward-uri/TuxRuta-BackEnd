package controller

import (
	"main/src/rutas/application"
	"main/src/rutas/domain/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ModifyRutaController struct {
	ModifyRutaUseCase *application.ModifyRutaUseCase
}

func NewModifyRutaController(modifyRutaUseCase *application.ModifyRutaUseCase) *ModifyRutaController {
	return &ModifyRutaController{
		ModifyRutaUseCase: modifyRutaUseCase,
	}
}

func (c *ModifyRutaController) HandleModifyRuta(g *gin.Context) {
	idparam := g.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid ruta ID"})
		return
	}

	var updatedRuta entities.Ruta
	if err := g.ShouldBindJSON(&updatedRuta); err != nil {
		g.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	err = c.ModifyRutaUseCase.Execute(id, &updatedRuta)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to modify ruta"})
		return
	}

	g.JSON(200, gin.H{"message": "Ruta modified successfully"})
}
