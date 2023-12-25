package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/johnstewart0820/jurassic_park/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type DinosaurController struct {
	DB *gorm.DB
}

func NewDinosaurController(DB *gorm.DB) DinosaurController {
	return DinosaurController{DB}
}

func (dc *DinosaurController) CreateDinosaur(ctx *gin.Context) {
	var payload *models.DinosaurRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var cage models.Cage
	cageResult := dc.DB.First(&cage, "id = ?", payload.CageId)
	if cageResult.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cage with that id exists"})
		return
	}

	var dinosaurs []models.Dinosaur
	result := dc.DB.Where(map[string]interface{}{"cage_id": cage.Id}).Find(&dinosaurs)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	if cage.Capacity < len(dinosaurs)+1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Cage is over limit."})
		return
	}

	if cage.Status == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Cage is down."})
		return
	}

	if !slices.Contains(models.DinosaurTypeList[payload.Type], payload.Spec) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Spec and Type of Dinosaur is wrong."})
		return
	}

	if payload.Type == "herbivore" {
		var dinosaurs []models.Dinosaur
		dc.DB.Where(map[string]interface{}{"type": "carnivore", "cage_id": payload.CageId}).Find(&dinosaurs)
		if len(dinosaurs) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "More than 1 carnivore is in cage"})
			return
		}
	} else if payload.Type == "carnivore" {
		var dinosaurs []models.Dinosaur
		dc.DB.Where(map[string]interface{}{"cage_id": payload.CageId, "type": "carnivore"}).Not(map[string]interface{}{"spec": []string{payload.Spec}}).Find(&dinosaurs)
		if len(dinosaurs) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Different spec of carnivore is in cage"})
			return
		}
		dc.DB.Where(map[string]interface{}{"cage_id": payload.CageId, "type": "herbivore"}).Find(&dinosaurs)
		if len(dinosaurs) > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "More than 1 herbivore is in cage"})
			return
		}
	}

	now := time.Now()
	newDinosaur := models.Dinosaur{
		Name:      payload.Name,
		Type:      payload.Type,
		Spec:      payload.Spec,
		CageId:    payload.CageId,
		Cage:      cage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result = dc.DB.Create(&newDinosaur)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newDinosaur})
}

func (dc *DinosaurController) FindDinosaurs(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	var spec = ctx.DefaultQuery("spec", "")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var dinosaurs []models.Dinosaur
	var results *gorm.DB
	if spec == "" {
		results = dc.DB.Limit(intLimit).Offset(offset).Find(&dinosaurs)
	} else {
		results = dc.DB.Where(map[string]interface{}{"spec": spec}).Limit(intLimit).Offset(offset).Find(&dinosaurs)
	}

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(dinosaurs), "data": dinosaurs})
}

func (dc *DinosaurController) FindDinosaurById(ctx *gin.Context) {
	dinosaurId := ctx.Param("id")

	var dinosaur models.Dinosaur
	result := dc.DB.First(&dinosaur, "id = ?", dinosaurId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No dinosaur with that id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dinosaur})
}

func (dc *DinosaurController) DeleteDinosaur(ctx *gin.Context) {
	dinosaurId := ctx.Param("id")

	result := dc.DB.Delete(&models.Dinosaur{}, "id = ?", dinosaurId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No dinosaur with that id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Deleted dinosaur successfully"})
}
