package routes

import (
	"main/src/user/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine, createUser *controllers.CreateUserHadler) {
	router.POST("/user", createUser.HandleCreateUser)
}
