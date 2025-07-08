package main

import (
	"log"
	"main/src/user/infraestructure"
	"main/src/user/infraestructure/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	infraestructure.InitDependeciesUser()

	router := gin.Default()
	routes.SetRoutes(router, infraestructure.CreateUserHandler,
		infraestructure.DeleteUserHandler, infraestructure.GetUsersHandler,
		infraestructure.GetUserByIDHandler, infraestructure.LoginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 5. Iniciar servidor
	log.Printf("🚀 Server started at :%s", port)
	log.Fatal(router.Run(":" + port))
}
