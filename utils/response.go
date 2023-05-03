package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response is the struct for an API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorData is the struct for error data i.e when Status is "error"
type ErrorData struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIResponse is a helper function to return an API response
func APIResponse(ctx *gin.Context, httpCode int, status string, message string, data interface{}) {
	ctx.JSON(httpCode, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// GetErrorMsg returns a list of meaningful error messages from binding tags.
// Reference: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Should be at least " + fe.Param() + " characters"
	case "max":
		return "Should be at most " + fe.Param() + " characters"
	}
	return "Unknown error"
}

// GetErrorData returns a list of error data
func GetErrorData(err error) []ErrorData {
	var errorData []ErrorData
	for _, fe := range err.(validator.ValidationErrors) {
		errorData = append(errorData, ErrorData{
			Field:   fe.Field(),
			Message: GetErrorMsg(fe),
		})
	}
	return errorData
}
