package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	db "github.com/paycrest/paycrest-protocol/sender/database"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/sender/ent"
	"github.com/paycrest/paycrest-protocol/sender/ent/enttest"
	"github.com/paycrest/paycrest-protocol/sender/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Set up a test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Create a new Gin router
	router := gin.New()

	// Set up mock dependencies
	ctrl := &UserController{}

	// Set up a test request with a JSON payload
	payload := &ent.User{
		Age:  25,
		Name: "John Doe",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Set up a test response recorder
	res := httptest.NewRecorder()

	// Call the CreateUser function
	router.POST("/users", ctrl.CreateUser)
	router.ServeHTTP(res, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, res.Code)
	var response utils.Response
	err = json.Unmarshal(res.Body.Bytes(), &response)

	// Assert the response body
	assert.NoError(t, err)
	assert.Equal(t, "User returned successfully", response.Message)
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "response.Data is not of type map[string]interface{}")
	assert.NotNil(t, data, "response.Data is nil")

	// Assert the response data
	assert.Contains(t, data, "id")
	assert.Equal(t, float64(payload.Age), data["age"].(float64))
	assert.Equal(t, payload.Name, data["name"].(string))
}
