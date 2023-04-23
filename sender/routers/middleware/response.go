package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// Response represents the custom JSON response structure
type APIResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// APIMiddleware is a middleware that sets the custom JSON response structure
func APIResponseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create a custom writer to capture the response body
		w := &responseWriter{body: gin.H{}}

		// Set the custom writer as the writer for the context
		ctx.Writer = w

		// Continue processing the request
		ctx.Next()

		// Retrieve the response data from the custom writer
		response := APIResponse{
			Status:  ctx.Writer.Status(),
			Data:    w.body,
			Message: "Success",
		}

		// Send the custom response as JSON
		ctx.JSON(ctx.Writer.Status(), response)
	}
}

// responseWriter is a custom writer that captures the response body
type responseWriter struct {
	gin.ResponseWriter
	body interface{}
}

// WriteJSON is a custom implementation of the WriteJSON method to capture the response body
func (w *responseWriter) WriteJSON(v interface{}) error {
	// Marshal the JSON data into bytes
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Set the response content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the bytes to the response
	_, err = w.ResponseWriter.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

