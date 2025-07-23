package routes

import (
	"main/src/infraestrucutre/middleware"
	"main/src/rutas/infraestructure/controller"
	"main/src/user/application/services"

	"github.com/gin-gonic/gin"
)

func SetRutaRoutes(
	router *gin.Engine,
	createRuta *controller.CreateRutaHandler,
	deleteRuta *controller.DeleteRutaController,
	getRutas *controller.GetRutasHandler,
	getRutaByID *controller.GetRutaByIDHandler,
	getRutaByName *controller.GetRutaByNombreHandler,
	modifyRuta *controller.ModifyRutaController,
	jwtService services.JWTService,

) {
	auth := middleware.AuthMiddleware(jwtService)
	router.POST("/rutas", auth, createRuta.HandleCreateRuta)
	router.DELETE("/rutas/:id", auth, deleteRuta.HandleDeleteRuta)
	router.GET("/rutas", auth, getRutas.HandleGetRutas)
	router.GET("/rutas/:id", auth, getRutaByID.HandleGetRutaByID)
	router.GET("/rutas/nombre/:nombre", auth, getRutaByName.HandleGetRutaByNombre)
	router.PUT("/rutas/:id", auth, modifyRuta.HandleModifyRuta)
}
