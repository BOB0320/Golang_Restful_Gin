package routes

import (
	"github.com/johnstewart0820/jurassic_park/controllers"

	"github.com/gin-gonic/gin"
)

type DinosaurRouteController struct {
	dinosaurController controllers.DinosaurController
}

func NewRoutedinosaurController(dinosaurController controllers.DinosaurController) DinosaurRouteController {
	return DinosaurRouteController{dinosaurController}
}

func (pc *DinosaurRouteController) SetupRoute(rg *gin.RouterGroup) {

	router := rg.Group("dinosaurs")
	router.POST("/", pc.dinosaurController.CreateDinosaur)
	router.GET("/", pc.dinosaurController.FindDinosaurs)
	router.GET("/:id", pc.dinosaurController.FindDinosaurById)
	router.DELETE("/:id", pc.dinosaurController.DeleteDinosaur)

}
