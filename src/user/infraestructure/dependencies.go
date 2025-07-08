package infraestructure

import (
	"main/src/core"
	"main/src/user/application"
	"main/src/user/infraestructure/controllers"
	"main/src/user/infraestructure/services"
)

var (
	CreateUserHandler  *controllers.CreateUserHadler
	DeleteUserHandler  *controllers.DeleteUserHandler
	GetUsersHandler    *controllers.GetUsersHandler
	GetUserByIDHandler *controllers.GetUserByIDHandler
	LoginHandler       *controllers.LoginHandler // ✅ Nuevo
)

func InitDependeciesUser() {
	core.InitPostgres()
	db := core.GetDB()

	// Servicios
	passwordService := services.NewBcryptPasswordService()
	jwtService := services.NewJWTService()
	// Repositorio
	userRepository := NewUserPostgreSQL(db)

	// Casos de uso existentes
	createUserUseCase := application.NewCreateUserUseCase(userRepository, passwordService)
	CreateUserHandler = controllers.NewCreateUserHandler(createUserUseCase)

	deleteUserUseCase := application.NewDeleteUserUseCase(userRepository)
	DeleteUserHandler = controllers.NewDeleteUserHandler(deleteUserUseCase)

	getAllUserUseCase := application.NewGetAllUserUseCase(userRepository)
	GetUsersHandler = controllers.NewGetUsersHandler(getAllUserUseCase)

	getUserByIdUseCase := application.NewGetUserByIdUseCase(userRepository)
	GetUserByIDHandler = controllers.NewGetUserByIDHandler(getUserByIdUseCase)

	loginUseCase := application.NewLoginUseCase(userRepository, passwordService, jwtService)
	LoginHandler = controllers.NewLoginHandler(loginUseCase)
}
