package infraestructure

import (
	"main/src/core"
	"main/src/paradas/application"
	"main/src/paradas/infraestructure/controller"
	rutasInfraestructure "main/src/rutas/infraestructure"
	"main/src/user/application/services"
	infraestructure_services "main/src/user/infraestructure/services"
)

var (
	CreateParadaHandler    *controller.CreateParadaHandler
	DeleteParadaHandler    *controller.DeleteParadaHandler
	GetParadasHandler      *controller.GetParadaHandler
	GetParadaByRutaHandler *controller.GetParadasAndRutasHandler
	ModifyParadaHandler    *controller.UpdateParadaHandler
	JWTService             services.JWTService
)

func InitDependeciesParada() {
	core.InitPostgres()
	db := core.GetDB()

	paradaRepository := NewParadasPostgreSQL(db)
	rutaRepository := rutasInfraestructure.NewRutasPostgreSQL(db)
	// Casos de uso
	createParadaUseCase := application.NewCreateParadaUseCase(paradaRepository, rutaRepository)

	deleteParadaUseCase := application.NewDeleteParadaUseCase(paradaRepository)

	getParadasUseCase := application.NewGetParadaUseCase(paradaRepository)

	getParadaByRutaUseCase := application.NewGetParadasAndRutaUseCase(paradaRepository, rutaRepository)

	modifyParadaUseCase := application.NewUpdateParadaUseCase(paradaRepository)

	// Handlers
	CreateParadaHandler = controller.NewCreateParadaHandler(createParadaUseCase)
	DeleteParadaHandler = controller.NewDeleteParadaHandler(deleteParadaUseCase)
	GetParadasHandler = controller.NewGetParadaHandler(getParadasUseCase)
	GetParadaByRutaHandler = controller.NewGetParadasAndRutasHandler(getParadaByRutaUseCase)
	ModifyParadaHandler = controller.NewUpdateParadaHandler(modifyParadaUseCase)

	JWTService = infraestructure_services.NewJWTService()
}
