package routes

import (
	"main/src/infraestrucutre/middleware"
	"main/src/paradas/infraestructure/controller"
	"main/src/user/application/services"

	"github.com/gin-gonic/gin"
)

func SetParadaRoutes(
	router *gin.Engine,
	createParada *controller.CreateParadaHandler,
	deleteParada *controller.DeleteParadaHandler,
	getParadas *controller.GetParadaHandler,
	getParadaByRuta *controller.GetParadasAndRutasHandler,
	modifyParada *controller.UpdateParadaHandler,
	jwtService services.JWTService,

) {

	auth := middleware.AuthMiddleware(jwtService)

	router.POST("/paradas", auth, createParada.HandleCreateParada)
	router.DELETE("/paradas/:id", auth, deleteParada.HandleDeleteParada)
	router.GET("/paradas/:id", getParadas.HandleGetParada)
	router.GET("/paradas/ruta/:ruta_id", getParadaByRuta.HandleGetParadasAndRutas)
	router.PUT("/paradas/:id", auth, modifyParada.HandleUpdateParada)
}
