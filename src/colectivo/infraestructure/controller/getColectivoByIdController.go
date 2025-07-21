package controller

import (
	"main/src/colectivo/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetColectivoByIDHandler struct {
	GetColectivoByIDUseCase *application.GetColectivoByIdUseCase
}

func NewGetColectivoByIDHandler(getColectivoByIDUseCase *application.GetColectivoByIdUseCase) *GetColectivoByIDHandler {
	return &GetColectivoByIDHandler{
		GetColectivoByIDUseCase: getColectivoByIDUseCase,
	}
}

func (c *GetColectivoByIDHandler) HandleGetColectivoByID(g *gin.Context) {
	idparam := g.Param("id")
	id, err := strconv.Atoi(idparam)
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid colectivo ID"})
		return
	}
	colectivo, err := c.GetColectivoByIDUseCase.Execute(id)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve colectivo"})
		return
	}
	g.JSON(200, colectivo)
}
