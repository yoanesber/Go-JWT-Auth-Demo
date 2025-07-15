package headers

import "github.com/gin-gonic/gin"

/**
* SecurityHeaders is a middleware function that sets various security-related HTTP headers
* to enhance the security of the web application.
* These headers help protect against common web vulnerabilities such as clickjacking, MIME type sniffing,
* cross-site scripting (XSS), and enforce secure connections.
 */
const (
	// Security headers
	xFrameOptions           = "X-Frame-Options"
	xContentTypeOptions     = "X-Content-Type-Options"
	xssProtection           = "X-XSS-Protection"
	strictTransportSecurity = "Strict-Transport-Security"
	referrerPolicy          = "Referrer-Policy"
	permissionsPolicy       = "Permissions-Policy"

	// Default values for security headers
	xFrameOptionsValue           = "DENY"
	xContentTypeOptionsValue     = "nosniff"
	xssProtectionValue           = "1; mode=block"
	strictTransportSecurityValue = "max-age=31536000; includeSubDomains; preload"
	referrerPolicyValue          = "no-referrer"
	permissionsPolicyValue       = "geolocation=(self), microphone=()"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// to protect against clickjacking attacks
		c.Writer.Header().Set(xFrameOptions, xFrameOptionsValue)

		// to prevent MIME type sniffing
		c.Writer.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)

		// to enable cross-site scripting (XSS) protection
		c.Writer.Header().Set(xssProtection, xssProtectionValue)

		// to enforce secure connections and
		// to ensure that browsers only connect to the server over HTTPS
		// This header is particularly important for production environments.
		c.Writer.Header().Set(strictTransportSecurity, strictTransportSecurityValue)

		// to control the referrer information sent with requests
		c.Writer.Header().Set(referrerPolicy, referrerPolicyValue)

		// to control which features can be used in the browser
		c.Writer.Header().Set(permissionsPolicy, permissionsPolicyValue)

		c.Next()
	}
}
