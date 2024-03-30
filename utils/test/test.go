package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// PerformRequest performs a http request with the given method, path, and payload
func PerformRequest(t *testing.T, method string, path string, payload interface{}, headers map[string]string, router *gin.Engine) (*httptest.ResponseRecorder, error) {
	req, _ := getRequest(method, path, payload)

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res, nil
}

// getRequest returns a new http.Request with the given method, path, and payload
func getRequest(method string, path string, payload interface{}) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
