package routes

import (
	"github.com/johnstewart0820/jurassic_park/controllers"

	"github.com/gin-gonic/gin"
)

type CageRouteController struct {
	cageController controllers.CageController
}

func NewCageRouteController(cageController controllers.CageController) CageRouteController {
	return CageRouteController{cageController}
}

func (rc *CageRouteController) SetupRoute(rg *gin.RouterGroup) {
	router := rg.Group("cages")

	router.GET("/", rc.cageController.FindCages)
	router.GET("/:id", rc.cageController.FindCageById)
	router.GET("/:id/dinosaurs", rc.cageController.FindDinosaursByCageId)
	router.POST("/", rc.cageController.CreateCage)
	router.PUT("/:id", rc.cageController.UpdateCage)
	router.DELETE("/:id", rc.cageController.DeleteCage)

}
