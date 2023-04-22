package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-services/infra/database"
	"github.com/paycrest/paycrest-services/models"
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
