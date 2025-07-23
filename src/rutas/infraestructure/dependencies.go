package infraestructure

import (
	"main/src/core"
	"main/src/rutas/application"
	"main/src/rutas/infraestructure/controller"
	"main/src/user/application/services"
	infraestructure_services "main/src/user/infraestructure/services"
)

var (
	CreateRutaHandler    *controller.CreateRutaHandler
	DeleteRutaHandler    *controller.DeleteRutaController
	GetRutasHandler      *controller.GetRutasHandler
	GetRutaByIDHandler   *controller.GetRutaByIDHandler
	ModifyRutaHandler    *controller.ModifyRutaController
	GetRutaByNameHandler *controller.GetRutaByNombreHandler
	JWTService           services.JWTService
)

func InitDependeciesRuta() {
	core.InitPostgres()
	db := core.GetDB()

	// Repositorio
	rutaRepository := NewRutasPostgreSQL(db)

	// Casos de uso
	createRutaUseCase := application.NewCreateRutaUseCase(rutaRepository)
	CreateRutaHandler = controller.NewCreateRutaHandler(createRutaUseCase)

	deleteRutaUseCase := application.NewDeleteRutaUseCase(rutaRepository)
	DeleteRutaHandler = controller.NewDeleteRutaHandler(deleteRutaUseCase)

	getRutasUseCase := application.NewGetRutasUseCase(rutaRepository)
	GetRutasHandler = controller.NewGetRutasHandler(getRutasUseCase)

	getRutaByIDUseCase := application.NewGetRutasByIdUseCase(rutaRepository)
	GetRutaByIDHandler = controller.NewGetRutaByIDHandler(getRutaByIDUseCase)

	modifyRutaUseCase := application.NewModifyRutaUseCase(rutaRepository)
	ModifyRutaHandler = controller.NewModifyRutaController(modifyRutaUseCase)

	getRutaByNombreUseCase := application.NewGetRutaByNombreUseCase(rutaRepository)
	GetRutaByNameHandler = controller.NewGetRutaByNombreHandler(getRutaByNombreUseCase)
	// Servicios
	JWTService = infraestructure_services.NewJWTService()
}
