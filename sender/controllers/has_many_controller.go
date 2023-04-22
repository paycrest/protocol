package controllers

import (
	"net/http"
	"strconv"

	"github.com/akmamun/gorm-pagination/pagination"
	"github.com/paycrest/paycrest-services/infra/database"

	"github.com/gin-gonic/gin"
	examples "github.com/paycrest/paycrest-services/models"
)

type CreditCardData struct {
	Number string `json:"number"`
}

//GetHasManyRelationUserData fetch user data with preload
func (ctrl *ExampleController) GetHasManyRelationUserData(ctx *gin.Context) {
	var user []examples.User
	// ctx.JSON(http.StatusOK, &user)
	// db :=base.DB.Preload("CreditCards").Find(&user)
	limit, _ := strconv.Atoi(ctx.GetString("limit"))
	offset, _ := strconv.Atoi(ctx.GetString("offset"))

	paginate := pagination.Paginate(&pagination.Param{
		DB:     database.DB,
		Limit:  int64(limit),
		Offset: int64(offset),
	}, &user)

	ctx.JSON(http.StatusOK, &paginate)

}

//GetHasManyRelationCreditCardData fetch credit-card data with preload
func (ctrl *ExampleController) GetHasManyRelationCreditCardData(ctx *gin.Context) {
	var creditCards []examples.CreditCard
	database.DB.Find(&creditCards)
	ctx.JSON(http.StatusOK, &creditCards)

}

// GetUserDetails based on user_id
func (ctrl *ExampleController) GetUserDetails(ctx *gin.Context) {
	var user []examples.User
	userId, _ := strconv.Atoi(ctx.DefaultQuery("user_id", ""))
	err := database.DB.Preload("CreditCards").First(&user, userId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": "Enter valid user"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
