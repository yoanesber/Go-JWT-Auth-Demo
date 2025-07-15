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
