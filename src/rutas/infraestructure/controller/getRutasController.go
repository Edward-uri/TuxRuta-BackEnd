package controller

import (
	"context"
	"main/src/core"
	"main/src/rutas/application"
	"time"

	"github.com/gin-gonic/gin"
)

type GetRutasHandler struct {
	GetRutasUseCase *application.GetRutasUseCase
}

func NewGetRutasHandler(getRutasUseCase *application.GetRutasUseCase) *GetRutasHandler {
	return &GetRutasHandler{
		GetRutasUseCase: getRutasUseCase,
	}
}

func (c *GetRutasHandler) HandleGetRutas(g *gin.Context) {
	cache := core.GetCacheService()
	cacheKey := "rutas_all"
	
	if cachedRutas, found := cache.Get(cacheKey); found {
		g.JSON(200, cachedRutas)
		return
	}

	ctx, cancel := context.WithTimeout(g.Request.Context(), 5*time.Second)
	defer cancel()

	rutas, err := c.GetRutasUseCase.Execute(ctx)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to retrieve rutas"})
		return
	}
	if len(rutas) == 0 {
		g.JSON(404, gin.H{"message": "No rutas found"})
		return
	}

	cache.Set(cacheKey, rutas, 1*time.Minute)

	g.JSON(200, rutas)
}
