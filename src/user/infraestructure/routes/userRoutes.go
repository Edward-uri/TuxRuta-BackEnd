package routes

import (
	"main/src/user/infraestructure/controller"

	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine, createUser *controller.CreateUserHadler,
	deleteUser *controller.DeleteUserHandler,
	getUsers *controller.GetUsersHandler,
	getUserByID *controller.GetUserByIDHandler,
	login *controller.LoginHandler) {
	router.POST("/user", createUser.HandleCreateUser)
	router.DELETE("/user/:id", deleteUser.HandleDeleteUser)
	router.GET("/users", getUsers.HandleGetUsers)
	router.GET("/user/:id", getUserByID.HandleGetUserByID)
	router.POST("/login", login.HandleLogin)
}
