package controller

import (
	"main/src/colectivo/application"
	"main/src/colectivo/domain/entities"

	"github.com/gin-gonic/gin"
)

type CreateColectivoHandler struct {
	CreateColectivoUseCase *application.CreateColectivoUseCase
}

func NewCreateColectivoHandler(createColectivoUseCase *application.CreateColectivoUseCase) *CreateColectivoHandler {
	return &CreateColectivoHandler{
		CreateColectivoUseCase: createColectivoUseCase,
	}
}

func (c *CreateColectivoHandler) HandleCreateColectivo(g *gin.Context) {
	var colectivo entities.Colectivo
	if err := g.ShouldBindJSON(&colectivo); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}

	creadoPor, exists := g.Get("user_id")
	if !exists {
		g.JSON(401, gin.H{"error": "Unauthorized"})
		print("Unauthorized access attempt", g.ClientIP())
		return
	}

	userID, ok := creadoPor.(int)
	if !ok {
		g.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	err := c.CreateColectivoUseCase.Execute(&colectivo, userID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(201, gin.H{"message": "Colectivo created successfully", "colectivo": colectivo})
}
