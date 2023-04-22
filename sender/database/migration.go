package database

import (
	"github.com/paycrest/paycrest-services/models"
)

//Add list of model add for migrations
//var migrationModels = []interface{}{&ex_models.Example{}, &model.Example{}, &model.Address{})}
var migrationModels = []interface{}{&models.Example{}, &models.User{}, &models.CreditCard{}}
