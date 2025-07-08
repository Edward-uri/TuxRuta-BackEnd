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
)

func InitDependeciesUser() {
	core.InitPostgres()
	db := core.GetDB()

	passwordService := services.NewBcryptPasswordService()

	userRepository := NewUserPostgreSQL(db)

	createUserUseCase := application.NewCreateUserUseCase(userRepository, passwordService)
	CreateUserHandler = controllers.NewCreateUserHandler(createUserUseCase)

	deleteUserUseCase := application.NewDeleteUserUseCase(userRepository)
	DeleteUserHandler = controllers.NewDeleteUserHandler(deleteUserUseCase)

	getAllUserUseCase := application.NewGetAllUserUseCase(userRepository)
	GetUsersHandler = controllers.NewGetUsersHandler(getAllUserUseCase)

	getUserByIdUseCase := application.NewGetUserByIdUseCase(userRepository)
	GetUserByIDHandler = controllers.NewGetUserByIDHandler(getUserByIdUseCase)
}
