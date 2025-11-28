package controller

import (
	"main/src/core"
	"main/src/paradas/application"
	"net/http"
	"time"

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
	// 1. Intentar obtener de caché
	cache := core.GetCacheService()
	cacheKey := "paradas_all"

	if cachedParadas, found := cache.Get(cacheKey); found {
		c.JSON(http.StatusOK, cachedParadas)
		return
	}

	paradas, err := h.GetParadasAllUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener paradas"})
		return
	}

	// 2. Guardar en caché por 1 minuto
	cache.Set(cacheKey, paradas, 1*time.Minute)

	c.JSON(http.StatusOK, paradas)
}
