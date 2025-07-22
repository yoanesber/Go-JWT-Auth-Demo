package http_util

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/logger"
)

// ErrorResponse represents the structure of an error response.
type HttpResponse struct {
	Message   string    `json:"message"`   // A user-friendly error message
	Error     any       `json:"error"`     // The actual error message (optional)
	Path      string    `json:"path"`      // The request path that caused the error (optional)
	Status    int       `json:"status"`    // HTTP status code (optional)
	Data      any       `json:"data"`      // Additional data related to the error (optional)
	Timestamp time.Time `json:"timestamp"` // The timestamp when the error occurred (optional)
}

/***** Basic Responses *****/
// Created sends a successful response with a 201 Created status.
// It is typically used when a new resource has been successfully created.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, HttpResponse{
		Message:   message,
		Error:     nil,
		Path:      c.Request.URL.Path,
		Status:    http.StatusCreated,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// Success sends a successful response with a 200 OK status.
// It is typically used for successful GET requests or other successful operations.
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, HttpResponse{
		Message:   message,
		Error:     nil,
		Path:      c.Request.URL.Path,
		Status:    http.StatusOK,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// BadRequest sends a 400 Bad Request response.
// It is typically used when the request cannot be processed due to client error.
func BadRequest(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusBadRequest, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusBadRequest,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// NotFound sends a 404 Not Found response.
// It is typically used when the requested resource cannot be found.
func NotFound(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusNotFound, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusNotFound,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// InternalServerError sends a 500 Internal Server Error response.
// It is typically used when an unexpected error occurs on the server.
func InternalServerError(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusInternalServerError, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusInternalServerError,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// Unauthorized sends a 401 Unauthorized response.
// It is typically used when authentication is required but has failed or has not been provided.
func Unauthorized(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusUnauthorized, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusUnauthorized,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// Forbidden sends a 403 Forbidden response.
// It is typically used when the server understands the request but refuses to authorize it.
func Forbidden(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusForbidden, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusForbidden,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// UnsupportedMediaType sends a 415 Unsupported Media Type response.
// It is typically used when the server refuses to accept the request because the payload format is invalid.
func UnsupportedMediaType(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusUnsupportedMediaType, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusUnsupportedMediaType,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// MethodNotAllowed sends a 405 Method Not Allowed response.
// It is typically used when the HTTP method used in the request is not allowed for the requested resource.
func MethodNotAllowed(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusMethodNotAllowed, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusMethodNotAllowed,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// Conflict sends a 409 Conflict response.
// It is typically used when a request could not be completed due to a conflict with the current state of the resource.
func Conflict(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusConflict, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusConflict,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// TooManyRequests sends a 429 Too Many Requests response.
// It is typically used when the user has sent too many requests in a given amount of time.
func TooManyRequests(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusTooManyRequests, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusTooManyRequests,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

// NoContent sends a 204 No Content response.
// It is typically used when the server successfully processes the request but does not need to return any content.
func NoContent(c *gin.Context, message string, err string) {
	logger.Error(err, nil)

	c.JSON(http.StatusNoContent, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusNoContent,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

/***** Map Responses *****/
func BadRequestMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Bad Request Map Error", nil)

	c.JSON(http.StatusBadRequest, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusBadRequest,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func NotFoundMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Not Found Map Error", nil)

	c.JSON(http.StatusNotFound, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusNotFound,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func InternalServerErrorMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Internal Server Error Map Error", nil)

	c.JSON(http.StatusInternalServerError, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusInternalServerError,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func UnauthorizedMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Unauthorized Map Error", nil)

	c.JSON(http.StatusUnauthorized, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusUnauthorized,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func ForbiddenMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Forbidden Map Error", nil)

	c.JSON(http.StatusForbidden, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusForbidden,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func UnsupportedMediaTypeMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Unsupported Media Type Map Error", nil)

	c.JSON(http.StatusUnsupportedMediaType, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusUnsupportedMediaType,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func MethodNotAllowedMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Method Not Allowed Map Error", nil)

	c.JSON(http.StatusMethodNotAllowed, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusMethodNotAllowed,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func ConflictMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Conflict Map Error", nil)

	c.JSON(http.StatusConflict, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusConflict,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func TooManyRequestsMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("Too Many Requests Map Error", nil)

	c.JSON(http.StatusTooManyRequests, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusTooManyRequests,
		Data:      nil,
		Timestamp: time.Now(),
	})
}

func NoContentMap(c *gin.Context, message string, err []map[string]string) {
	logger.Error("No Content Map Error", nil)

	c.JSON(http.StatusNoContent, HttpResponse{
		Message:   message,
		Error:     err,
		Path:      c.Request.URL.Path,
		Status:    http.StatusNoContent,
		Data:      nil,
		Timestamp: time.Now(),
	})
}
