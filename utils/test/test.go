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
func PerformRequest(t *testing.T, method string, path string, payload interface{}, auth *string, router *gin.Engine) (*httptest.ResponseRecorder, error) {
	req, _ := GetRequest(t, method, path, payload, router)

	if auth != nil {
		req.Header.Set("Authorization", "Bearer "+*auth)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res, nil
}

// GetRequest returns a new http.Request with the given method, path, and payload
func GetRequest(t *testing.T, method string, path string, payload interface{}, router *gin.Engine) (*http.Request, error) {
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
