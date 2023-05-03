package accounts

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	db "github.com/paycrest/paycrest-protocol/database"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Set up a test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Set up test router
	router := gin.New()
	ctrl := &AuthController{}
	router.POST("/register", ctrl.Register)

	t.Run("with valid payload", func(t *testing.T) {
		// Test register with valid payload
		payload := RegisterPayload{
			FirstName: "Ike",
			LastName:  "Ayo",
			Email:     "ikeayo@example.com",
			Password:  "password",
		}

		res, err := test.PerformRequest(t, "POST", "/register", payload, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusCreated, res.Code)

		var response utils.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User created successfully", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is not of type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		// Assert the response data
		assert.Contains(t, data, "id")
		match, _ := regexp.MatchString(
			`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
			data["id"].(string),
		)
		if !match {
			t.Errorf("Expected '%s' to be a valid UUID", data["id"].(string))
		}
		assert.Equal(t, payload.Email, data["email"].(string))
		assert.Equal(t, payload.FirstName, data["first_name"].(string))
		assert.Equal(t, payload.LastName, data["last_name"].(string))
	})

	t.Run("with existing user", func(t *testing.T) {
		// Test register with existing user
		payload := RegisterPayload{
			FirstName: "Ike",
			LastName:  "Ayo",
			Email:     "ikeayo@example.com",
			Password:  "password",
		}

		res, err := test.PerformRequest(t, "POST", "/register", payload, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusBadRequest, res.Code)

		var response utils.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User with email already exists", response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("with invalid email", func(t *testing.T) {
		// Test register with invalid email
		payload := RegisterPayload{
			FirstName: "Ike",
			LastName:  "Ayo",
			Email:     "invalid-email",
			Password:  "password",
		}

		res, err := test.PerformRequest(t, "POST", "/register", payload, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusBadRequest, res.Code)

		var response utils.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to validate payload", response.Message)
		assert.Equal(t, "error", response.Status)
		data, ok := response.Data.([]interface{})
		assert.True(t, ok, "response.Data is not of type []interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		// Assert the response errors in data
		assert.Len(t, data, 1)
		errorMap, ok := data[0].(map[string]interface{})
		assert.True(t, ok, "error is not of type map[string]interface{}")
		assert.NotNil(t, errorMap, "error is nil")
		assert.Contains(t, errorMap, "field")
		assert.Equal(t, "Email", errorMap["field"].(string))
		assert.Contains(t, errorMap, "message")
		assert.Equal(t, "Must be a valid email address", errorMap["message"].(string))
	})

}
