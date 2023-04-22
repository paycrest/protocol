package controllers

import (
	"github.com/akmamun/gin-boilerplate-examples/infra/database"
	"github.com/akmamun/gin-boilerplate-examples/models"
	"github.com/akmamun/gorm-pagination/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ctrl *ExampleController) GetExamplePaginated(ctx *gin.Context) {
	var example []models.Example

	limit, _ := strconv.Atoi(ctx.GetString("limit"))
	offset, _ := strconv.Atoi(ctx.GetString("offset"))

	paginateData := pagination.Paginate(&pagination.Param{
		DB:      database.DB,
		Offset:  int64(offset),
		Limit:   int64(limit),
		OrderBy: "id desc",
	}, &example)

	ctx.JSON(http.StatusOK, paginateData)

}
