package authorization

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	metacontext "github.com/yoanesber/go-consumer-api-with-jwt/pkg/context-data/meta-context"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
	jwtutil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/jwt-util"
)

/**
* JwtValidation is a middleware function that validates JWT tokens in the request header.
* It checks if the token is present, has the correct format, and is valid.
* If the token is valid, it extracts user information from the token claims and injects it into the request context.
* If the token is invalid or missing, it returns an unauthorized error response.
 */
var (
	TokenType string
	JWTSecret string
)

// LoadEnv loads environment variables
func LoadEnv() {
	TokenType = os.Getenv("TOKEN_TYPE")
	JWTSecret = os.Getenv("JWT_SECRET")
}

func JwtValidation() gin.HandlerFunc {
	// Load environment variables
	LoadEnv()

	return func(c *gin.Context) {
		// Get the token from the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httputil.Unauthorized(c, "No token provided", "Authorization header is missing")
			c.Abort()
			return
		}

		// Check if the token starts with TokenType
		tokenPrefix := TokenType + " "
		if !strings.HasPrefix(authHeader, tokenPrefix) {
			httputil.Unauthorized(c, "Invalid token format", fmt.Sprintf("Token must start with '%s'", tokenPrefix))
			c.Abort()
			return
		}

		// Extract the token string
		tokenStr := strings.TrimPrefix(authHeader, tokenPrefix)
		if tokenStr == "" {
			httputil.Unauthorized(c, "Invalid token format", "Token string is empty")
			c.Abort()
			return
		}

		// Parse the token and validate it
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// For HS256 signing method
			if token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
				// Validate the token signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Return the secret key for validation
				return []byte(JWTSecret), nil
			}

			// For RS256 signing method
			// Load the public key from the environment variable
			publicKey, err := jwtutil.LoadPublicKey()
			if err != nil {
				return nil, err
			}

			// Validate the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the public key for validation
			return publicKey, nil
		})

		if err != nil {
			httputil.Unauthorized(c, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Check if the token is valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			httputil.Unauthorized(c, "Invalid token", "Token is not valid")
			c.Abort()
			return
		}

		// Get the user ID from the claims
		// Convert the user ID to int64
		userID, _ := jwtutil.GetInt64Claim(claims, "userid")

		// Inject user information into the request context
		meta := metacontext.UserInformationMeta{
			UserID:   userID,
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			Roles:    jwtutil.GetStringSliceClaim(claims, "roles"),
		}
		ctx := metacontext.InjectUserInformationMeta(c.Request.Context(), meta)

		// Set the new request context with user information
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
