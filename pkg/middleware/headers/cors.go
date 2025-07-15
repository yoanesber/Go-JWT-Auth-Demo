package headers

import "github.com/gin-gonic/gin"

/**
* CorsHeaders is a middleware that sets Cross-Origin Resource Sharing (CORS) headers
* to allow cross-origin requests from the frontend (e.g., from a different domain or port).
* It is typically used in web applications to enable communication between the frontend and backend
* when they are hosted on different origins (domains, protocols, or ports).
 */
const (
	// CORS headers
	accessControlAllowOrigin      = "Access-Control-Allow-Origin"
	accessControlMaxAge           = "Access-Control-Max-Age"
	accessControlAllowMethods     = "Access-Control-Allow-Methods"
	accessControlAllowHeaders     = "Access-Control-Allow-Headers"
	accessControlExposeHeaders    = "Access-Control-Expose-Headers"
	accessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// Default values for CORS headers
	accessControlAllowOriginValue      = "http://localhost"
	accessControlMaxAgeValue           = "86400" // 1 day in seconds
	accessControlAllowMethodsValue     = "POST, GET, OPTIONS, PUT, DELETE, UPDATE"
	accessControlAllowHeadersValue     = "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token"
	accessControlExposeHeadersValue    = "Content-Length"
	accessControlAllowCredentialsValue = "true"
)

func CorsHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(accessControlAllowOrigin, accessControlAllowOriginValue)
		c.Writer.Header().Set(accessControlMaxAge, accessControlMaxAgeValue)
		c.Writer.Header().Set(accessControlAllowMethods, accessControlAllowMethodsValue)
		c.Writer.Header().Set(accessControlAllowHeaders, accessControlAllowHeadersValue)
		c.Writer.Header().Set(accessControlExposeHeaders, accessControlExposeHeadersValue)
		c.Writer.Header().Set(accessControlAllowCredentials, accessControlAllowCredentialsValue)

		if c.Request.Method == "OPTIONS" {
			// Handle preflight request
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}
