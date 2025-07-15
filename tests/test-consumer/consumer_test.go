package test_consumer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/handler"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/service"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/middleware/authorization"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
)

const (
	dummyAdminToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ5b3VyX2p3dF9hdWRpZW5jZSIsImVtYWlsIjoiYWRtaW5AbXlnbWFpbC5jb20iLCJleHAiOjE3NTA2NTAzNjEsImlhdCI6MTc1MDQ3NzU2MSwiaXNzIjoieW91cl9qd3RfaXNzdWVyIiwicm9sZXMiOlsiUk9MRV9BRE1JTiJdLCJzdWIiOiJhZG1pbiIsInVzZXJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiJ9.iBUMUUbwUy2CswqmR23hCNBF872cLjcn12UrUWJEm34"
	dummyNonAdminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ5b3VyX2p3dF9hdWRpZW5jZSIsImVtYWlsIjoidXNlcm9uZUBteWdtYWlsLmNvbSIsImV4cCI6MTc1MDY1MDMyOSwiaWF0IjoxNzUwNDc3NTI5LCJpc3MiOiJ5b3VyX2p3dF9pc3N1ZXIiLCJyb2xlcyI6WyJST0xFX1VTRVIiXSwic3ViIjoidXNlcm9uZSIsInVzZXJpZCI6MiwidXNlcm5hbWUiOiJ1c2Vyb25lIn0.1ZA8dS7Eb5Hn4PaZagTsSesqwGt_tplXLntW9QPVYeo"
	dummyInvalidToken  = "invalid.token.string"
	dummyEmptyToken    = ""
	dummyExpiredToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ5b3VyX2p3dF9hdWRpZW5jZSIsImVtYWlsIjoidXNlcm9uZUBteWdtYWlsLmNvbSIsImV4cCI6MTc1MDQ3NzUyOSwiaWF0IjoxNzUwNDc3NTI5LCJpc3MiOiJ5b3VyX2p3dF9pc3N1ZXIiLCJyb2xlcyI6WyJST0xFX1VTRVIiXSwic3ViIjoidXNlcm9uZSIsInVzZXJpZCI6MiwidXNlcm5hbWUiOiJ1c2Vyb25lIn0.V3DfjAgw7kNCBP1ueidv9lJV5s4J491hSDERWj3hlKE"
)

func TestGetAllConsumers_Success(t *testing.T) {
	// Define a mocked repository, service, and handler
	// This will allow us to test the handler without needing a real database connection
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint with the JWT token in the Authorization header
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	req.Header.Set("Authorization", "Bearer "+dummyAdminToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, httpResponse.Data)
	assert.Nil(t, httpResponse.Error)
}

func TestGetAllConsumers_Unauthorized(t *testing.T) {
	// Define a mocked repository, service, and handler
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint without a token
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestGetAllConsumers_Forbidden(t *testing.T) {
	// Define a mocked repository, service, and handler
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint with a non-admin token
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	req.Header.Set("Authorization", "Bearer "+dummyNonAdminToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestGetAllConsumers_InvalidToken(t *testing.T) {
	// Define a mocked repository, service, and handler
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint with an invalid token
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	req.Header.Set("Authorization", "Bearer "+dummyInvalidToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestGetAllConsumers_EmptyToken(t *testing.T) {
	// Define a mocked repository, service, and handler
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint with an empty token
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	req.Header.Set("Authorization", "Bearer "+dummyEmptyToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}

func TestGetAllConsumers_ExpiredToken(t *testing.T) {
	// Define a mocked repository, service, and handler
	r := NewConsumerMockedRepository()
	s := service.NewConsumerService(r)
	h := handler.NewConsumerHandler(s)

	// Set up the Gin router and the route for getting all consumers
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(authorization.JwtValidation())
	router.GET("/api/v1/consumers", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.GetAllConsumers)

	// Create a request to the endpoint with an expired token
	req, _ := http.NewRequest("GET", "/api/v1/consumers", nil)
	req.Header.Set("Authorization", "Bearer "+dummyExpiredToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Unmarshal the response body into a HttpResponse struct
	var httpResponse httputil.HttpResponse
	err := json.Unmarshal(w.Body.Bytes(), &httpResponse)
	assert.NoError(t, err)
	assert.Empty(t, httpResponse.Data)
	assert.NotNil(t, httpResponse.Error)
}
