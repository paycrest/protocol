package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/paycrest/protocol/types"
)

const (
	HTTP_RETRY_ATTEMPTS = 3
	HTTP_RETRY_INTERVAL = 5
)

// APIResponse is a helper function to return an API response
func APIResponse(ctx *gin.Context, httpCode int, status string, message string, data interface{}) {
	ctx.JSON(httpCode, types.Response{
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
	case "oneof":
		options := strings.Split(fe.Param(), ",")
		return "Must be one of " + strings.Join(options, ", ")
	}
	return "Unknown error"
}

// GetErrorData returns a list of error data
func GetErrorData(err error) []types.ErrorData {
	var errorData []types.ErrorData
	for _, fe := range err.(validator.ValidationErrors) {
		errorData = append(errorData, types.ErrorData{
			Field:   fe.Field(),
			Message: GetErrorMsg(fe),
		})
	}
	return errorData
}

// MakeJSONRequest makes a JSON request
func MakeJSONRequest(ctx context.Context, method, url string, payload map[string]interface{}, headers map[string]string) (responseData map[string]interface{}, err error) {
	if !ContainsString([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}, method) {
		return nil, errors.New("invalid method")
	}

	// Create a new request
	requestBody, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new context and add the request to it
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	var res *http.Response
	for i := 0; i < HTTP_RETRY_ATTEMPTS; i++ { // On failure, retry up to 3 times
		res, err = client.Do(req)
		if err == nil && res.StatusCode < 500 && res.StatusCode != 429 {
			break
		}
		if i < HTTP_RETRY_ATTEMPTS-1 { // Avoid sleep after the last attempt
			time.Sleep(HTTP_RETRY_INTERVAL * time.Second) // Wait for 5 seconds before retrying
		}
	}
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 500 { // Return on server errors
		return nil, fmt.Errorf(fmt.Sprintf("server error: %d", res.StatusCode))
	}
	if res.StatusCode >= 400 { // Return on client errors
		return nil, fmt.Errorf(fmt.Sprintf("client error: %d", res.StatusCode))
	}

	// Decode the response body into a map
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(responseBody, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Paginate parses the pagination query params and returns the offset(page) and limit(pageSize)
func Paginate(ctx *gin.Context) (page int, pageSize int) {
	// Parse pagination query params
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err2 := strconv.Atoi(ctx.Query("pageSize"))

	// Set defaults if not provided
	if err != nil || page < 1 {
		page = 1
	}
	if err2 != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calculate offsets
	page = (page - 1) * pageSize

	return page, pageSize
}

// IsURL checks if a string is a valid URL
func IsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}
