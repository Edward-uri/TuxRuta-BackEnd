package main

import (
	"log"
	"os"

	colectivo_infraestructure "main/src/colectivo/infraestructure"
	colectivoRoutes "main/src/colectivo/infraestructure/routes"
	rutasInfraestructure "main/src/rutas/infraestructure"
	rutasRoutes "main/src/rutas/infraestructure/routes"
	"main/src/user/infraestructure"
	userRoutes "main/src/user/infraestructure/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	infraestructure.InitDependeciesUser()
	colectivo_infraestructure.InitDependeciesColectivo()
	rutasInfraestructure.InitDependeciesRuta()
	router := gin.Default()

	userRoutes.SetRoutes(
		router,
		infraestructure.CreateUserHandler,
		infraestructure.DeleteUserHandler,
		infraestructure.GetUsersHandler,
		infraestructure.GetUserByIDHandler,
		infraestructure.LoginHandler,
	)

	colectivoRoutes.SetColectivoRoutes(
		router,
		colectivo_infraestructure.CreateColectivoHandler,
		colectivo_infraestructure.DeleteColectivoHandler,
		colectivo_infraestructure.GetColectivosHandler,
		colectivo_infraestructure.GetColectivoByIDHandler,
		colectivo_infraestructure.GetColectivoByMatriculaHandler,
		colectivo_infraestructure.ModifyColectivoHandler,
		colectivo_infraestructure.JWTService,
	)

	rutasRoutes.SetRutaRoutes(
		router,
		rutasInfraestructure.CreateRutaHandler,
		rutasInfraestructure.DeleteRutaHandler,
		rutasInfraestructure.GetRutasHandler,
		rutasInfraestructure.GetRutaByIDHandler,
		rutasInfraestructure.GetRutaByNameHandler,
		rutasInfraestructure.ModifyRutaHandler,
		rutasInfraestructure.JWTService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server started at :%s", port)
	log.Fatal(router.Run(":" + port))
}
