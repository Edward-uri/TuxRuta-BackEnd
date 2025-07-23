package controller

import (
	"main/src/paradas/application"
	"main/src/paradas/domain/entities"

	"github.com/gin-gonic/gin"
)

type UpdateParadaHandler struct {
	UpdateParadaUseCase *application.UpdateParadaUseCase
}

func NewUpdateParadaHandler(updateParadaUseCase *application.UpdateParadaUseCase) *UpdateParadaHandler {
	return &UpdateParadaHandler{
		UpdateParadaUseCase: updateParadaUseCase,
	}
}

func (c *UpdateParadaHandler) HandleUpdateParada(g *gin.Context) {
	var request struct {
		ID        int                `json:"id"`
		Nombre    string             `json:"nombre"`
		Ubicacion entities.Ubicacion `json:"ubicacion"`
		RutaID    int                `json:"ruta_id"`
	}

	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	parada := entities.Parada{
		ID:        request.ID,
		Nombre:    request.Nombre,
		Ubicacion: request.Ubicacion,
		RutaID:    request.RutaID,
	}

	err := c.UpdateParadaUseCase.Execute(parada)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, gin.H{"message": "Parada updated successfully"})
}
