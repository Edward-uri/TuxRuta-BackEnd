package main

import (
	"log"
	"os"

	colectivo_infraestructure "main/src/colectivo/infraestructure"
	colectivoRoutes "main/src/colectivo/infraestructure/routes"
	paradasInfraestructure "main/src/paradas/infraestructure"
	paradasRoutes "main/src/paradas/infraestructure/routes"
	rutasInfraestructure "main/src/rutas/infraestructure"
	rutasRoutes "main/src/rutas/infraestructure/routes"
	"main/src/user/infraestructure"
	userRoutes "main/src/user/infraestructure/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:4200",
			"http://3.228.56.73:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Middleware de compresión GZIP
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	infraestructure.InitDependeciesUser()
	colectivo_infraestructure.InitDependeciesColectivo()
	rutasInfraestructure.InitDependeciesRuta()
	paradasInfraestructure.InitDependeciesParada()

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

	paradasRoutes.SetParadaRoutes(
		router,
		paradasInfraestructure.CreateParadaHandler,
		paradasInfraestructure.DeleteParadaHandler,
		paradasInfraestructure.GetParadasHandler,
		paradasInfraestructure.GetParadaByRutaHandler,
		paradasInfraestructure.ModifyParadaHandler,
		paradasInfraestructure.GetParadasAllHandler,
		paradasInfraestructure.JWTService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started at :%s", port)
	log.Fatal(router.Run(":" + port))
}
