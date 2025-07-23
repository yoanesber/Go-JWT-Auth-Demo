package headers

import (
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	httputil "github.com/yoanesber/go-jwt-auth-demo/pkg/util/http-util"
)

/**
* CorsHeaders is a middleware that sets Cross-Origin Resource Sharing (CORS) headers
* to allow cross-origin requests from the frontend (e.g., from a different domain or port).
* It is typically used in web applications to enable communication between the frontend and backend
* when they are hosted on different origins (domains, protocols, or ports).
 */

func CorsHeaders() gin.HandlerFunc {
	env := os.Getenv("NODE_ENV")

	var allowedOrigins []string
	if env == "production" {
		allowedOrigins = strings.Split(os.Getenv("FRONTEND_URL_PRODUCTION"), ",")
	} else {
		allowedOrigins = strings.Split(os.Getenv("FRONTEND_URL"), ",")
	}

	// Set CORS headers for allowed origins
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			httputil.BadRequest(c, "Missing Origin", "The request does not have an Origin header")
			c.Abort()
			return
		}

		// Validate the origin URL
		// Ensure the origin is a valid URL and uses HTTP or HTTPS scheme
		requestedOrigin, err := url.Parse(origin)
		if err != nil {
			httputil.BadRequest(c, "Invalid Origin", "The provided origin is not a valid URL")
			c.Abort()
			return
		}
		if requestedOrigin.Scheme != "http" && requestedOrigin.Scheme != "https" {
			httputil.BadRequest(c, "Invalid Origin", "The provided origin must use HTTP or HTTPS scheme")
			c.Abort()
			return
		}

		// Check if the origin is in the allowed origins list
		// If the origin is allowed, set CORS headers
		for _, allowed := range allowedOrigins {
			allowedOriginTrimmed := strings.TrimSpace(allowed)
			if origin == allowedOriginTrimmed || allowedOriginTrimmed == "*" {
				maxAge := 24 * time.Hour
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Access-Control-Max-Age", maxAge.String())

				if c.Request.Method == "OPTIONS" {
					httputil.NoContent(c, "Preflight request successful", "CORS preflight request handled successfully")
					c.Abort()
					return
				}

				c.Next()
				return
			}
		}

		// If the origin is not allowed, respond with an error
		httputil.Forbidden(c, "CORS Error", "Origin not allowed")
		c.Abort()
	}

}
