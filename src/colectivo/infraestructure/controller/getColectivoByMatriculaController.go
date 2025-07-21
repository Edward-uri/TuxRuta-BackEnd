package controller

import (
	"main/src/colectivo/application"

	"github.com/gin-gonic/gin"
)

type GetColectivoByMatriculaHandler struct {
	GetColectivoByMatriculaUseCase *application.GetColectivoByMatriculaUseCase
}

func NewGetColectivoByMatriculaHandler(getColectivoByMatriculaUseCase *application.GetColectivoByMatriculaUseCase) *GetColectivoByMatriculaHandler {
	return &GetColectivoByMatriculaHandler{
		GetColectivoByMatriculaUseCase: getColectivoByMatriculaUseCase,
	}
}

func (c *GetColectivoByMatriculaHandler) HandleGetColectivoByMatricula(g *gin.Context) {
	matricula := g.Param("matricula")
	if matricula == "" {
		g.JSON(400, gin.H{"error": "Matricula is required"})
		return
	}
	colectivo, err := c.GetColectivoByMatriculaUseCase.Execute(matricula)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve colectivo"})
		return
	}
	g.JSON(200, colectivo)
}
