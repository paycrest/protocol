package controllers

import (
	"github.com/akmamun/gin-boilerplate-examples/infra/database"
	"github.com/akmamun/gin-boilerplate-examples/models"
	"github.com/gin-gonic/gin"
	"net/http"
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
