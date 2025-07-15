package authorization

import (
	"github.com/gin-gonic/gin"

	metacontext "github.com/yoanesber/go-consumer-api-with-jwt/pkg/context-data/meta-context"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
)

/**
* RoleBasedAccessControl is a middleware function that checks if the user has the required roles to access a resource.
* It extracts user metadata from the request context and checks if the user has any of the allowed roles.
* If the user has at least one of the allowed roles, it allows the request to proceed to the next handler.
* If the user does not have any of the allowed roles, it returns a forbidden response and aborts the request.
 */
func RoleBasedAccessControl(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no allowed roles are provided, allow access
		if len(allowedRoles) == 0 {
			c.Next()
			return
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractUserInformationMeta(c.Request.Context())
		if !ok {
			httputil.InternalServerError(c, "Failed to extract metadata", "Unable to extract user metadata from context")
			c.Abort()
			return
		}

		// Get the user roles from the metadata
		userRoles := meta.Roles
		if len(userRoles) == 0 {
			httputil.Forbidden(c, "No roles found", "User does not have any roles")
			c.Abort()
			return
		}

		// Check if the user has any of the allowed roles
		// If the user has at least one allowed role, proceed to the next handler
		for _, role := range userRoles {
			for _, allowed := range allowedRoles {
				if role == allowed {
					c.Next()
					return
				}
			}
		}

		// If the user does not have any of the allowed roles, return a forbidden response
		// and abort the request
		httputil.Forbidden(c, "Access denied", "User does not have the required role")
		c.Abort()
	}
}
