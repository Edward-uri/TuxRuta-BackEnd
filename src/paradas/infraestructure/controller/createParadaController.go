package controller

import (
	"main/src/paradas/application"
	"main/src/paradas/domain/entities"

	"github.com/gin-gonic/gin"
)

type CreateParadaHandler struct {
	CreateParadaUseCase *application.CreateParadaUseCase
}

func NewCreateParadaHandler(createParadaUseCase *application.CreateParadaUseCase) *CreateParadaHandler {
	return &CreateParadaHandler{
		CreateParadaUseCase: createParadaUseCase,
	}
}

func (c *CreateParadaHandler) HandleCreateParada(g *gin.Context) {
	var request struct {
		Nombre    string             `json:"nombre"`
		Ubicacion entities.Ubicacion `json:"ubicacion"`
		RutaID    int                `json:"ruta_id"`
	}

	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	creadoPor, exists := g.Get("user_id")
	if !exists {
		g.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := creadoPor.(int)
	if !ok {
		g.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	parada, err := c.CreateParadaUseCase.Execute(request.Nombre, request.Ubicacion, request.RutaID, userID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(201, gin.H{"message": "Parada created successfully", "parada": parada})
}
