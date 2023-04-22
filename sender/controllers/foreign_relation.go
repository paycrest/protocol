package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-services/infra/database"
	"github.com/paycrest/paycrest-services/models"
)

// GetNormalData get normal data if added pagination see example_controller
func (ctrl *ExampleController) GetNormalData(ctx *gin.Context) {
	var categories []models.Category
	db := database.GetDB()
	db.Find(&categories)
	ctx.JSON(http.StatusOK, gin.H{"data": categories})

}

// GetForeignRelationData Get Foreign Data with Preload
func (ctrl *ExampleController) GetForeignRelationData(ctx *gin.Context) {
	var articles []models.Article
	db := database.GetDB()
	db.Preload("Category").Find(&articles)
	ctx.JSON(http.StatusOK, &articles)

}
