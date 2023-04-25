package controllers

import (
	"log"
	"net/http"

	db "github.com/paycrest/paycrest-protocol/sender/database"
	"github.com/paycrest/paycrest-protocol/sender/ent"
	utils "github.com/paycrest/paycrest-protocol/sender/utils"
	"github.com/paycrest/paycrest-protocol/sender/utils/logger"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	var payload ent.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		utils.APIResponse(ctx, http.StatusBadRequest, 400, "error", err.Error())
		return
	}

	user, err := db.Client.User.
		Create().
		SetAge(payload.Age).
		SetName(payload.Name).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		utils.APIResponse(ctx, http.StatusInternalServerError, 500, "error", err.Error())
		return
	}

	utils.APIResponse(ctx, http.StatusOK, 200, "User returned successfully", &user)
}

func (ctrl *UserController) GetUsers(ctx *gin.Context) {
	var users []*ent.User

	var err error
	users, err = db.Client.User.Query().All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		utils.APIResponse(ctx, http.StatusInternalServerError, 500, "error", err.Error())
		return
	}
	log.Println("users", users)
	utils.APIResponse(ctx, http.StatusOK, 200, "Users returned successfully", &users)
}
