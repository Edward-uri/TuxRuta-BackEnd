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
	// 1. Cargar variables de entorno PRIMERO
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// 2. Inicializar dependencias ANTES de crear las rutas
	infraestructure.InitDependeciesUser()

	// 3. Crear router e inicializar rutas DESPUÉS
	router := gin.Default()
	routes.SetRoutes(router, infraestructure.CreateUserHandler)

	// 4. Obtener puerto del environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 5. Iniciar servidor
	log.Printf("🚀 Server started at :%s", port)
	log.Fatal(router.Run(":" + port))
}
