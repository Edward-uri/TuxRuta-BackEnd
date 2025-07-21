package controller

import (
	"main/src/colectivo/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteColectivoHandler struct {
	DeleteColectivoUseCase *application.DeleteColectivoUseCase
}

func NewDeleteColectivoHandler(deleteColectivoUseCase *application.DeleteColectivoUseCase) *DeleteColectivoHandler {
	return &DeleteColectivoHandler{
		DeleteColectivoUseCase: deleteColectivoUseCase,
	}
}

func (d *DeleteColectivoHandler) HandleDeleteColectivo(g *gin.Context) {
	idparam := g.Param("id")
	id, err := strconv.Atoi(idparam)
	if id == 0 || err != nil {
		g.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	err = d.DeleteColectivoUseCase.Execute(id)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"message": "Colectivo deleted successfully"})
}
