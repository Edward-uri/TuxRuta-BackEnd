package infraestructure

import (
	"main/src/colectivo/application"
	"main/src/colectivo/infraestructure/controller"
	controllers "main/src/colectivo/infraestructure/controller"
	"main/src/core"
	"main/src/user/application/services"
	infraestructure_services "main/src/user/infraestructure/services"
)

var (
	CreateColectivoHandler         *controller.CreateColectivoHandler
	DeleteColectivoHandler         *controller.DeleteColectivoHandler
	GetColectivosHandler           *controller.GetColectivosHandler
	GetColectivoByIDHandler        *controller.GetColectivoByIDHandler
	ModifyColectivoHandler         *controller.ModifyColectivoController
	GetColectivoByMatriculaHandler *controller.GetColectivoByMatriculaHandler
	JWTService                     services.JWTService
)

func InitDependeciesColectivo() {
	core.InitPostgres()
	db := core.GetDB()

	// Repositorio
	colectivoRepository := NewColectivoPostgreSQL(db)

	// Casos de uso
	createColectivoUseCase := application.NewCreateColectivoUseCase(colectivoRepository)
	CreateColectivoHandler = controllers.NewCreateColectivoHandler(createColectivoUseCase)

	deleteColectivoUseCase := application.NewDeleteColectivoUseCase(colectivoRepository)
	DeleteColectivoHandler = controllers.NewDeleteColectivoHandler(deleteColectivoUseCase)

	getAllColectivosUseCase := application.NewGetColectivosUseCase(colectivoRepository)
	GetColectivosHandler = controllers.NewGetColectivosHandler(getAllColectivosUseCase)

	getColectivoByIDUseCase := application.NewGetColectivoByIdUseCase(colectivoRepository)
	GetColectivoByIDHandler = controllers.NewGetColectivoByIDHandler(getColectivoByIDUseCase)

	getColectivoByMatriculaUseCase := application.NewGetColectivoByMatriculaUseCase(colectivoRepository)
	GetColectivoByMatriculaHandler = controllers.NewGetColectivoByMatriculaHandler(getColectivoByMatriculaUseCase)

	modifyColectivoUseCase := application.NewModifyColectivoUseCase(colectivoRepository)
	ModifyColectivoHandler = controllers.NewModifyColectivoController(modifyColectivoUseCase)

	// Inicializa el servicio JWT (usa tu constructor)
	JWTService = infraestructure_services.NewJWTService()
}
