package headers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
)

/**
* SecurityHeaders is a middleware function that sets various security-related HTTP headers
* to enhance the security of the web application.
* These headers help protect against common web vulnerabilities such as clickjacking, MIME type sniffing,
* cross-site scripting (XSS), and enforce secure connections.
 */

func SecurityHeaders() gin.HandlerFunc {
	isSSLRedirect := os.Getenv("IS_SSL") == "TRUE"

	secureMiddleware := secure.New(secure.Options{
		// Protects against reflected XSS attacks in older browsers
		BrowserXssFilter: true,

		// Prevents MIME sniffing by browsers
		ContentTypeNosniff: true,

		// Prevents the site from being framed to mitigate clickjacking attacks
		FrameDeny: true,

		// Redirect all HTTP traffic to HTTPS (enabled in production only)
		// Enable only in production
		SSLRedirect: isSSLRedirect,

		// Recognize HTTPS requests when behind a reverse proxy like Nginx
		// Required if using a proxy (safe in all environments)
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},

		// Enables HTTP Strict Transport Security (HSTS) for one year
		// Enable only in production with HTTPS
		STSSeconds: 31536000,

		// Applies HSTS to all subdomains
		// Enable only in production with HTTPS and full subdomain coverage
		STSIncludeSubdomains: true,

		// Preload HSTS configuration into browsers
		// Enable **only if** your domain meets preload requirements (https://hstspreload.org/)
		STSPreload: true,

		// Always send HSTS header, even over HTTP
		// Set to true only if not using a reverse proxy that terminates TLS
		ForceSTSHeader: false,

		// Controls how much referrer information is included with requests
		// "no-referrer" is most private, use "strict-origin-when-cross-origin" for more flexibility
		ReferrerPolicy: "no-referrer",

		// Controls access to browser features (geolocation, mic, camera, etc.)
		// Safe and recommended; adjust per app requirements
		PermissionsPolicy: "geolocation=(self), microphone=(), camera=()",

		// Restricts which resources can be loaded and embedded
		// Highly recommended for frontend or API returning HTML
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'",

		// Ensures top-level document doesn't share context group with cross-origin documents
		// Recommended for modern browser security
		CrossOriginOpenerPolicy: "same-origin",

		// Prevents embedding cross-origin resources unless explicitly allowed
		// Recommended for performance and security
		CrossOriginEmbedderPolicy: "require-corp",

		// Restricts which origins can load your resources (images, scripts, etc.)
		// Recommended for stricter resource isolation
		CrossOriginResourcePolicy: "same-origin",

		// Disables DNS prefetching to reduce privacy leaks
		// Safe in most cases; can impact performance slightly
		XDNSPrefetchControl: "off",

		// Blocks Flash/Adobe cross-domain policies (legacy)
		// Safe to disable in modern web apps
		XPermittedCrossDomainPolicies: "none",
	})

	return func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			httputil.BadRequest(c, "Security Header Error", err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
