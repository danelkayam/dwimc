package integration

import (
	"bytes"
	api_model "dwimc/internal/api/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type NoValidatedResponse struct {
	Data  any                      `json:"data"`
	Error *api_model.ErrorResponse `json:"error"`
}

func PerformOKRequest[T any](
	t *testing.T,
	router *gin.Engine,
	method string,
	url string,
	apiKey string,
	payload any,
) T {
	w := performRequest(router, method, url, apiKey, payload)
	assert.Equal(t, http.StatusOK, w.Code)

	var response api_model.Response[T]

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Failed parsing response body")
	assert.NotNilf(t, response.Data, "Response Data is nil")
	assert.Nilf(t, response.Error, "Response Error is not nil")

	return response.Data
}

func PerformOKRequestNoValidateResponse(
	t *testing.T,
	router *gin.Engine,
	method string,
	url string,
	apiKey string,
	payload any,
) *NoValidatedResponse {
	w := performRequest(router, method, url, apiKey, payload)
	assert.Equal(t, http.StatusOK, w.Code)

	var response NoValidatedResponse

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Failed parsing response body")

	return &response
}

func PerformFailedRequest(
	t *testing.T,
	router *gin.Engine,
	method string,
	url string,
	apiKey string,
	payload any,
	expectedStatusCode int,
) *api_model.ErrorResponse {
	w := performRequest(router, method, url, apiKey, payload)
	assert.Equal(t, expectedStatusCode, w.Code)

	var response api_model.Response[any]

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoErrorf(t, err, "Failed parsing response body")
	assert.Nilf(t, response.Data, "Response Data is not nil")
	assert.NotNilf(t, response.Error, "Response Error is nil")

	return response.Error
}

func performRequest(
	router *gin.Engine,
	method string,
	url string,
	apiKey string,
	payload any,
) *httptest.ResponseRecorder {
	body := []byte{}
	if payload != nil {
		body, _ = json.Marshal(payload)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("X-API-Key", apiKey)

	router.ServeHTTP(w, req)
	return w
}
