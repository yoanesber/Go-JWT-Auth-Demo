package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	metacontext "github.com/yoanesber/go-consumer-api-with-jwt/pkg/context-data/meta-context"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/logger"
)

/**
* RequestLogger is a middleware function that logs incoming HTTP requests.
* It initializes the logger, records the request details, and logs them after the request is processed.
 */
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process the request first
		// This allows the middleware to log the request details after the request has been processed
		// This is important to capture the response status and duration accurately
		c.Next()

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractUserInformationMeta(c.Request.Context())
		if !ok {
			meta = metacontext.UserInformationMeta{
				Username: "unknown",
				Roles:    []string{"unknown"},
			}
		}

		// Get the username from the context
		// This assumes that the username is set in the context by JWT validation middleware
		if meta.Username == "" {
			meta.Username = "unknown"
		}

		// Then log the request details
		// This is done after the request is processed to capture the response status and duration
		duration := time.Since(start)
		logger.RequestLogger.WithFields(logrus.Fields{
			"content_length": c.Request.ContentLength,
			"content_type":   c.ContentType(),
			"duration":       duration.String(),
			"ip":             c.ClientIP(),
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"query":          c.Request.URL.Query(),
			"referer":        c.Request.Referer(),
			"request_id":     c.Writer.Header().Get("X-Request-Id"),
			"status":         c.Writer.Status(),
			"user_agent":     c.Request.UserAgent(),
			"username":       meta.Username,
			"roles":          meta.Roles,
		}).Info("Incoming request")
	}
}
