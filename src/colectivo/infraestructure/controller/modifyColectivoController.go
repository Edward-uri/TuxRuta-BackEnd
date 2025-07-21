package controller

import (
	"main/src/colectivo/application"
	"main/src/colectivo/domain/entities"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ModifyColectivoController struct {
	ModifyColectivoUseCase *application.ModifyColectivoUseCase
}

func NewModifyColectivoController(modifyColectivoUseCase *application.ModifyColectivoUseCase) *ModifyColectivoController {
	return &ModifyColectivoController{
		ModifyColectivoUseCase: modifyColectivoUseCase,
	}
}

func (c *ModifyColectivoController) HandleModifyColectivo(g *gin.Context) {
	idparam := g.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid colectivo ID"})
		return
	}

	var colectivoData entities.Colectivo
	if err := g.ShouldBindJSON(&colectivoData); err != nil {
		g.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	err = c.ModifyColectivoUseCase.Execute(id, &colectivoData)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to modify colectivo"})
		return
	}

	g.JSON(200, gin.H{"message": "Colectivo modified successfully"})
}
