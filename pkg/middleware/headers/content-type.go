package headers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	httputil "github.com/yoanesber/go-jwt-auth-demo/pkg/util/http-util"
)

/**
 * ContentType is a middleware function that checks the Content-Type header of incoming requests.
 * It ensures that the Content-Type is set to `application/json` for POST, PUT, and PATCH requests.
 * If the Content-Type is not set correctly, it returns a 415 Unsupported Media Type error and aborts the request.
 * This middleware is useful for enforcing the expected content type for API requests.
 */

func ContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		contentType := c.GetHeader("Content-Type")

		// Only enforce for methods that require a body
		if method == http.MethodPost || method == http.MethodPut {
			if !strings.HasPrefix(contentType, "application/json") {
				httputil.UnsupportedMediaType(c, "Unsupported Media Type", "Content-Type must be `application/json`")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
