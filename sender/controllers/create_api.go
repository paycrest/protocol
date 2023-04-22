package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-services/infra/database"
	"github.com/paycrest/paycrest-services/models"
	"github.com/paycrest/paycrest-services/sender/utils/logger"
)

type ExampleController struct{}

func (ctrl *ExampleController) CreateExample(ctx *gin.Context) {
	example := new(models.Example)

	err := ctx.ShouldBindJSON(&example)
	if err != nil {
		logger.Errorf("error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = database.DB.Create(&example).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &example)
}

func (ctrl *ExampleController) GetExampleData(ctx *gin.Context) {
	var examples []models.Example
	database.DB.Find(&examples)
	ctx.JSON(http.StatusOK, gin.H{"data": examples})

}
