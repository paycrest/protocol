package controllers

import (
	"net/http"

	"github.com/paycrest/paycrest-services/sender/database"
	"github.com/paycrest/paycrest-services/sender/models"
	utils "github.com/paycrest/paycrest-services/sender/utils"
	"github.com/paycrest/paycrest-services/sender/utils/logger"

	"github.com/gin-gonic/gin"
)

type ExampleController struct{}

func (ctrl *ExampleController) CreateExample(ctx *gin.Context) {
	example := new(models.Example)

	err := ctx.ShouldBindJSON(&example)
	if err != nil {
		logger.Errorf("error: %v", err)
		utils.APIResponse(ctx, http.StatusInternalServerError, 500, "error", err.Error())
		return
	}
	err = database.DB.Create(&example).Error
	if err != nil {
		logger.Errorf("error: %v", err)
		utils.APIResponse(ctx, http.StatusInternalServerError, 500, "error", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &example)
}

func (ctrl *ExampleController) GetExampleData(ctx *gin.Context) {
	var examples []models.Example
	database.DB.Find(&examples)
	utils.APIResponse(ctx, http.StatusOK, 200, "Examples returned successfully", examples)
}
