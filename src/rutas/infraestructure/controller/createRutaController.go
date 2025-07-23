package controller

import (
	"main/src/rutas/application"
	"main/src/rutas/domain/entities"

	"github.com/gin-gonic/gin"
)

type CreateRutaHandler struct {
	CreateRutaUseCase *application.CreateRutaUseCase
}

func NewCreateRutaHandler(createRutaUseCase *application.CreateRutaUseCase) *CreateRutaHandler {
	return &CreateRutaHandler{
		CreateRutaUseCase: createRutaUseCase,
	}
}

func (c *CreateRutaHandler) HandleCreateRuta(g *gin.Context) {
	var request struct {
		Nombre      string           `json:"nombre"`
		Descripcion string           `json:"descripcion"`
		Points      []entities.Point `json:"points"`
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

	pathData := entities.PathData{Points: request.Points}

	ruta, err := c.CreateRutaUseCase.Execute(request.Nombre, request.Descripcion, pathData, userID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(201, gin.H{"message": "Ruta created successfully", "ruta": ruta})
}
