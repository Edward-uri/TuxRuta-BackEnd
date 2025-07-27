package controller

import (
	"main/src/paradas/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetParadasAllHandler struct {
	GetParadasAllUseCase *application.GetParadasAllUseCase
}

func NewGetParadasAllHandler(GetParadasAllUseCase *application.GetParadasAllUseCase) *GetParadasAllHandler {
	return &GetParadasAllHandler{
		GetParadasAllUseCase: GetParadasAllUseCase,
	}
}

func (h *GetParadasAllHandler) HandleGetParadas(c *gin.Context) {
	paradas, err := h.GetParadasAllUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener paradas"})
		return
	}
	c.JSON(http.StatusOK, paradas)
}
