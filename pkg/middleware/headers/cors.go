package headers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
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
		allowedOrigins = []string{os.Getenv("FRONTEND_URL_PRODUCTION")}
	} else {
		allowedOrigins = []string{os.Getenv("FRONTEND_URL")}
	}

	// Set CORS headers for allowed origins
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
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
