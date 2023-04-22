package controllers

import (
	"net/http"
	"strconv"

	"github.com/paycrest/paycrest-services/sender/database"
	"github.com/paycrest/paycrest-services/sender/models"

	"github.com/akmamun/gorm-pagination/pagination"
	"github.com/gin-gonic/gin"
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
