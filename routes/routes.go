package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/handler"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/service"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/middleware/authorization"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/middleware/headers"
	"github.com/yoanesber/go-consumer-api-with-jwt/pkg/middleware/logging"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
)

// SetupRouter initializes the router and sets up the routes for the application.
func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()

	// Set up middleware for the router
	// Middleware is used to handle cross-cutting concerns such as logging, security, and request ID generation
	r.Use(
		headers.SecurityHeaders(),
		headers.CorsHeaders(),
		headers.ContentType(),
		logging.RequestLogger(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	// Set up the authentication routes
	// These routes handle user login and authentication
	authGroup := r.Group("/auth")
	{
		// Routes for authentication
		// These routes handle user login
		s := service.NewAuthService()
		h := handler.NewAuthHandler(s)

		// Define the routes for authentication
		// These routes handle user login
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh-token", h.RefreshToken)
	}

	// Set up the API version 1 routes
	v1 := r.Group("/api/v1", authorization.JwtValidation())
	{
		// Routes for consumer management
		// These routes handle CRUD operations for consumers
		consumerGroup := v1.Group("/consumers")
		{
			// Initialize the transaction repository and service
			// This is where the actual implementation of the repository and service would be used
			r := repository.NewConsumerRepository()
			s := service.NewConsumerService(r)

			// Initialize the transaction handler with the service
			// This handler handles the HTTP requests and responses for transaction-related operations
			h := handler.NewConsumerHandler(s)

			// Define the routes for transaction management
			// These routes handle CRUD operations for transactions
			// The GET methods are accessible to both admin and user roles
			consumerGroup.GET("", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), h.GetAllConsumers)
			consumerGroup.GET("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), h.GetConsumerByID)
			consumerGroup.GET("/active", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), h.GetActiveConsumers)
			consumerGroup.GET("/inactive", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), h.GetInactiveConsumers)
			consumerGroup.GET("/suspended", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), h.GetSuspendedConsumers)

			// The POST and PUT methods are restricted to admin users only
			consumerGroup.POST("", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.CreateConsumer)
			consumerGroup.PATCH("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN"), h.UpdateConsumerStatus)
		}
	}

	// NoRoute handler for undefined routes
	// This handler will be called when no other route matches the request
	r.NoRoute(func(c *gin.Context) {
		httputil.NotFound(c, "Not Found", "The requested resource was not found")
	})

	// NoMethod handler for unsupported HTTP methods
	// This handler will be called when a request method is not allowed for the requested resource
	r.NoMethod(func(c *gin.Context) {
		httputil.MethodNotAllowed(c, "Method Not Allowed", "The requested method is not allowed for this resource")
	})

	return r
}
