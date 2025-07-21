package routes

import (
	"main/src/colectivo/infraestructure/controller"
	"main/src/infraestrucutre/middleware"
	"main/src/user/application/services" // importa la interfaz JWTService

	"github.com/gin-gonic/gin"
)

func SetColectivoRoutes(
	router *gin.Engine,
	createColectivo *controller.CreateColectivoHandler,
	deleteColectivo *controller.DeleteColectivoHandler,
	getColectivos *controller.GetColectivosHandler,
	getColectivoByID *controller.GetColectivoByIDHandler,
	getColectivoByMatricula *controller.GetColectivoByMatriculaHandler,
	modifyColectivo *controller.ModifyColectivoController,
	jwtService services.JWTService, // <-- usa la interfaz aquí
) {
	auth := middleware.AuthMiddleware(jwtService)

	router.POST("/colectivo", auth, createColectivo.HandleCreateColectivo)
	router.DELETE("/colectivo/:id", auth, deleteColectivo.HandleDeleteColectivo)
	router.PUT("/colectivo/:id", auth, modifyColectivo.HandleModifyColectivo)

	router.GET("/colectivos", getColectivos.HandleGetColectivos)
	router.GET("/colectivo/matricula/:matricula", getColectivoByMatricula.HandleGetColectivoByMatricula)
	router.GET("/colectivo/:id", getColectivoByID.HandleGetColectivoByID)
}
