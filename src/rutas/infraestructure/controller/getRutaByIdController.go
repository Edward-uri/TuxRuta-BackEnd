package controller

import (
	"main/src/rutas/application"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetRutaByIDHandler struct {
	useCase *application.GetRutasByIdUseCase
}

func NewGetRutaByIDHandler(useCase *application.GetRutasByIdUseCase) *GetRutaByIDHandler {
	return &GetRutaByIDHandler{useCase: useCase}
}

func (h *GetRutaByIDHandler) HandleGetRutaByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Invalid ruta ID"})
		return
	}

	ruta, err := h.useCase.Execute(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve ruta"})
		return
	}

	c.JSON(200, ruta)
}
