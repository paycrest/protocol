package controllers

import (
	"net/http"

	"github.com/paycrest/paycrest-services/sender/database"
	"github.com/paycrest/paycrest-services/sender/models"

	"github.com/gin-gonic/gin"
)

// SelectedFiledFetch fields fetch from defining new struct
type SelectedFiledFetch struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

func (ctrl *ExampleController) GetSelectedFieldData(ctx *gin.Context) {
	var selectData []SelectedFiledFetch
	database.DB.Model(&models.Article{}).Find(&selectData)
	ctx.JSON(http.StatusOK, selectData)

}
