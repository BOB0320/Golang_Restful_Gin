package main

import (
	"log"

	"github.com/johnstewart0820/jurassic_park/controllers"
	"github.com/johnstewart0820/jurassic_park/initializers"
	"github.com/johnstewart0820/jurassic_park/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine

	CageController      controllers.CageController
	CageRouteController routes.CageRouteController

	DinosaurController      controllers.DinosaurController
	DinosaurRouteController routes.DinosaurRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	CageController = controllers.NewCageController(initializers.DB)
	CageRouteController = routes.NewCageRouteController(CageController)

	DinosaurController = controllers.NewDinosaurController(initializers.DB)
	DinosaurRouteController = routes.NewRoutedinosaurController(DinosaurController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	routerV1 := server.Group("/api/v1")

	CageRouteController.SetupRoute(routerV1)
	DinosaurRouteController.SetupRoute(routerV1)
	log.Fatal(server.Run(":" + config.ServerPort))
}
