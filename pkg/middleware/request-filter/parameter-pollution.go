package request_filter

import (
	"fmt"

	"github.com/gin-gonic/gin"

	httputil "github.com/yoanesber/go-jwt-auth-demo/pkg/util/http-util"
)

/**
 * DetectParameterPollution is a middleware function that checks for parameter pollution in the request.
 * It looks for query parameters that have multiple values (e.g., ?param=value1&param=value2).
 * If it detects any parameters with multiple values, it considers it as parameter pollution and returns a 400 Bad Request response.
 * This middleware is useful for preventing parameter pollution attacks where an attacker might try to manipulate query parameters.
 */
func DetectParameterPollution() gin.HandlerFunc {
	return func(c *gin.Context) {
		pollutedParams := make(map[string][]string)

		query := c.Request.URL.Query()
		for key, values := range query {
			if len(values) > 1 {
				pollutedParams[key] = values
			}
		}

		if len(pollutedParams) > 0 {
			httputil.BadRequest(c, "Parameter Pollution Detected", fmt.Sprintf("Parameter pollution detected: %v", pollutedParams))
			c.Abort()
			return
		}

		c.Next()
	}
}
