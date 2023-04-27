package controllers

import (
	"log"
	"net/http"

	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	var payload ent.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "error", err.Error())
		return
	}

	user, err := db.Client.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetEmail(payload.Email).
		SetPassword(payload.Password).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "error", err.Error())
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "User returned successfully", &user)
}

func (ctrl *UserController) GetUsers(ctx *gin.Context) {
	var users []*ent.User

	var err error
	users, err = db.Client.User.Query().All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "error", err.Error())
		return
	}
	log.Println("users", users)
	u.APIResponse(ctx, http.StatusOK, "success", "Users returned successfully", &users)
}
