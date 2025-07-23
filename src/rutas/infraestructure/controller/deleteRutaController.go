package controller

import (
	"main/src/rutas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteRutaController struct {
	deleteRutaUseCase *application.DeleteRutaUseCase
}

func NewDeleteRutaHandler(deleteRutaUseCase *application.DeleteRutaUseCase) *DeleteRutaController {
	return &DeleteRutaController{deleteRutaUseCase: deleteRutaUseCase}
}

func (d *DeleteRutaController) HandleDeleteRuta(g *gin.Context) {
	idparam := g.Param("id")
	id, err := strconv.Atoi(idparam)
	if id == 0 || err != nil {
		g.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	if err := d.deleteRutaUseCase.Execute(id); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{"message": "Ruta deleted successfully"})
}
