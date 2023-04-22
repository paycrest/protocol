package controllers

import (
	"net/http"
	"strconv"

	"github.com/akmamun/gorm-pagination/pagination"
	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-services/infra/database"
	"github.com/paycrest/paycrest-services/models"
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
