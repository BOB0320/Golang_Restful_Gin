package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/johnstewart0820/jurassic_park/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CageController struct {
	DB *gorm.DB
}

func NewCageController(DB *gorm.DB) CageController {
	return CageController{DB}
}

func (cg *CageController) CreateCage(ctx *gin.Context) {
	var payload *models.CageRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newCage := models.Cage{
		Status:    payload.Status,
		Capacity:  payload.Capacity,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := cg.DB.Create(&newCage)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newCage})
}

func (cg *CageController) UpdateCage(ctx *gin.Context) {
	cageId := ctx.Param("id")

	var payload *models.CageRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var updatedCage models.Cage
	result := cg.DB.First(&updatedCage, "id = ?", cageId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cage with that id exists"})
		return
	}

	if payload.Status == false {
		var dinosaurs []models.Dinosaur
		cg.DB.Where(map[string]interface{}{"cage_id": cageId}).Find(&dinosaurs)
		if len(dinosaurs) > 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "The cage can not down."})
			return
		}
	}

	now := time.Now()
	fmt.Println(payload.Status)
	cageToUpdate := models.Cage{
		Status:    payload.Status,
		Capacity:  payload.Capacity,
		CreatedAt: updatedCage.CreatedAt,
		UpdatedAt: now,
	}

	result = cg.DB.Model(&updatedCage).Updates(cageToUpdate)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": cageToUpdate})
}

func (cg *CageController) FindCageById(ctx *gin.Context) {
	cageId := ctx.Param("id")

	var cage models.Cage
	result := cg.DB.First(&cage, "id = ?", cageId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cage with that id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": cage})
}

func (cg *CageController) FindDinosaursByCageId(ctx *gin.Context) {
	cageId := ctx.Param("id")

	var cage models.Cage
	result := cg.DB.First(&cage, "id = ?", cageId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cage with that id exists"})
		return
	}

	var dinosaurs []models.Dinosaur

	result = cg.DB.Where(map[string]interface{}{"cage_id": cageId}).Find(&dinosaurs)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dinosaurs})
}

func (cg *CageController) FindCages(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")
	var status = ctx.DefaultQuery("status", "")
	var parsedStatus bool
	if status == "true" {
		parsedStatus = true
	} else {
		parsedStatus = false
	}

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var cages []models.Cage
	var results *gorm.DB
	if status == "" {
		results = cg.DB.Limit(intLimit).Offset(offset).Find(&cages)
	} else {
		results = cg.DB.Where(map[string]interface{}{"status": parsedStatus}).Limit(intLimit).Offset(offset).Find(&cages)
	}

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(cages), "data": cages})
}

func (cg *CageController) DeleteCage(ctx *gin.Context) {
	cageId := ctx.Param("id")

	cg.DB.Delete(&models.Dinosaur{}, "cage_id = ?", cageId)

	result := cg.DB.Delete(&models.Cage{}, "id = ?", cageId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Deleted Cage and Dinosaurs in it successfully"})
}
