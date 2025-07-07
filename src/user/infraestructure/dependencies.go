package infraestructure

import (
	"main/src/core"
	"main/src/user/application"
	"main/src/user/infraestructure/controllers"
	"main/src/user/infraestructure/services"
)

var (
	CreateUserHandler *controllers.CreateUserHadler
)

func InitDependeciesUser() {
	core.InitPostgres()
	db := core.GetDB()

	passwordService := services.NewBcryptPasswordService()

	userRepository := NewUserPostgreSQL(db)

	createUserUseCase := application.NewCreateUserUseCase(userRepository, passwordService)

	CreateUserHandler = controllers.NewCreateUserHandler(createUserUseCase)
}
