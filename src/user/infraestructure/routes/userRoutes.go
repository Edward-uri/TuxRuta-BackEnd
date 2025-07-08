package routes

import (
	"main/src/user/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine, createUser *controllers.CreateUserHadler,
	deleteUser *controllers.DeleteUserHandler,
	getUsers *controllers.GetUsersHandler,
	getUserByID *controllers.GetUserByIDHandler,
	login *controllers.LoginHandler) {
	router.POST("/user", createUser.HandleCreateUser)
	router.DELETE("/user/:id", deleteUser.HandleDeleteUser)
	router.GET("/users", getUsers.HandleGetUsers)
	router.GET("/user/:id", getUserByID.HandleGetUserByID)
	router.POST("/login", login.HandleLogin)
}
