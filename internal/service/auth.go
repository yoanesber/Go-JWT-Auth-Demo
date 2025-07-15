package service

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/config/database"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/repository"
	jwtutil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/jwt-util"
)

var (
	once              sync.Once
	JWTSecret         string
	TokenType         string
	SigningMethod     string
	JWTAudience       string
	JWTIssuer         string
	JWTExpirationHour string
	AccessTokenTTL    time.Duration
)

// LoadEnv loads environment variables
func LoadEnv() {
	once.Do(func() {
		JWTSecret = os.Getenv("JWT_SECRET")
		TokenType = os.Getenv("TOKEN_TYPE")
		SigningMethod = os.Getenv("JWT_ALGORITHM")
		JWTAudience = os.Getenv("JWT_AUDIENCE")
		JWTIssuer = os.Getenv("JWT_ISSUER")
		JWTExpirationHour = os.Getenv("JWT_EXPIRATION_HOUR")

		// Load access and refresh token TTL from environment variables
		access, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))
		AccessTokenTTL = time.Duration(access) * time.Minute
	})
}

// Interface for auth service
// This interface defines the methods that the auth service should implement
type AuthService interface {
	Login(loginReq entity.LoginRequest) (entity.LoginResponse, error)
	RefreshToken(refreshTokenReq entity.RefreshTokenRequest) (entity.RefreshTokenResponse, error)
}

// This struct defines the AuthService that contains a user repository and a role repository
// It implements the AuthService interface and provides methods for authentication-related operations
type authService struct{}

// NewAuthService creates a new instance of AuthService with the given user and role repositories.
// It initializes the authService struct and returns it.
func NewAuthService() AuthService {
	return &authService{}
}

// Login authenticates a user with the given username and password.
// It retrieves the token for the user if the authentication is successful.
func (s *authService) Login(loginReq entity.LoginRequest) (entity.LoginResponse, error) {
	// Load environment variables
	LoadEnv()

	// Get the database connection from the context
	db := database.GetPostgres()
	if db == nil {
		return entity.LoginResponse{}, fmt.Errorf("database connection is nil")
	}

	// Validate the authentication parameters using the validation
	if err := loginReq.Validate(); err != nil {
		return entity.LoginResponse{}, err
	}

	var tokenStr string
	var refreshTokenStr string
	var expirationDateStr string
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		userRepo := repository.NewUserRepository()
		userService := NewUserService(userRepo)
		existingUser, err := userService.GetUserByUsername(loginReq.Username)
		if err != nil {
			return err
		}

		// Check some conditions for the user
		if existingUser.Equals(&entity.User{}) {
			return fmt.Errorf("user with username %s not found", loginReq.Username)
		}
		if !*existingUser.IsEnabled {
			return fmt.Errorf("user with username %s is not enabled", loginReq.Username)
		}
		if !*existingUser.IsAccountNonExpired {
			return fmt.Errorf("user account is expired")
		}
		if !*existingUser.IsAccountNonLocked {
			return fmt.Errorf("user account is locked")
		}
		if !*existingUser.IsCredentialsNonExpired {
			return fmt.Errorf("user credentials are expired")
		}
		if *existingUser.IsDeleted {
			return fmt.Errorf("user with username %s is deleted", loginReq.Username)
		}

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginReq.Password)); err != nil {
			return fmt.Errorf("invalid credentials for user %s", loginReq.Username)
		}

		// Generate an access token for the user
		tokenStr, err = GenerateJWTToken(existingUser)
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %w", err)
		}

		// Parse the JWT token
		jwtToken, err := ParseJWTToken(tokenStr)
		if err != nil {
			return fmt.Errorf("failed to parse JWT token: %w", err)
		}

		// Get the expiration date from the token
		expirationDateStr, err = GetExpirationDateFromToken(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to get expiration date from token: %w", err)
		}

		// Generate a refresh token for the user
		refreshTokenRepo := repository.NewRefreshTokenRepository()
		refreshTokenService := NewRefreshTokenService(refreshTokenRepo)
		jwtRefreshToken, err := refreshTokenService.CreateRefreshToken(existingUser.ID)
		if err != nil {
			return fmt.Errorf("failed to create refresh token: %w", err)
		}
		if jwtRefreshToken.Equals(&entity.RefreshToken{}) {
			return fmt.Errorf("failed to create refresh token")
		}

		refreshTokenStr = jwtRefreshToken.Token

		// Update the last login time for the user
		_, err = userService.UpdateLastLogin(existingUser.ID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to update last login time: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		AccessToken:    tokenStr,
		RefreshToken:   refreshTokenStr,
		ExpirationDate: expirationDateStr,
		TokenType:      TokenType,
	}, nil
}

// RefreshToken refreshes the access token using the provided refresh token.
// It retrieves the new access token and refresh token for the user.
func (s *authService) RefreshToken(refreshTokenReq entity.RefreshTokenRequest) (entity.RefreshTokenResponse, error) {
	// Load environment variables
	LoadEnv()

	// Get the database connection from the context
	db := database.GetPostgres()
	if db == nil {
		return entity.RefreshTokenResponse{}, fmt.Errorf("database connection is nil")
	}

	// Validate the refresh token request
	if err := refreshTokenReq.Validate(); err != nil {
		return entity.RefreshTokenResponse{}, err
	}

	var accessTokenStr string
	var refreshTokenStr string
	var expirationDateStr string
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the refresh token exists
		refreshTokenRepo := repository.NewRefreshTokenRepository()
		refreshTokenService := NewRefreshTokenService(refreshTokenRepo)
		existingRefreshToken, err := refreshTokenService.GetRefreshTokenByToken(refreshTokenReq.RefreshToken)
		if err != nil {
			return err
		}
		if existingRefreshToken.Equals(&entity.RefreshToken{}) {
			return fmt.Errorf("refresh token not found")
		}

		// If found, check if the refresh token is expired
		ok, _ := refreshTokenService.VerifyExpirationDate(existingRefreshToken.ExpiryDate)
		if !ok {
			return fmt.Errorf("refresh token is expired")
		}

		// Get user details using the user ID from the refresh token
		userRepo := repository.NewUserRepository()
		userService := NewUserService(userRepo)
		userDetails, err := userService.GetUserByID(existingRefreshToken.UserID)
		if err != nil {
			return err
		}
		if userDetails.Equals(&entity.User{}) {
			return fmt.Errorf("user with ID %d not found", existingRefreshToken.UserID)
		}

		// Generate an access token for the user
		accessTokenStr, err = GenerateJWTToken(userDetails)
		if err != nil {
			return fmt.Errorf("failed to generate JWT token: %w", err)
		}

		// Parse the JWT token
		jwtToken, err := ParseJWTToken(accessTokenStr)
		if err != nil {
			return fmt.Errorf("failed to parse JWT token: %w", err)
		}

		// Get the expiration date from the token
		expirationDateStr, err = GetExpirationDateFromToken(jwtToken)
		if err != nil {
			return fmt.Errorf("failed to get expiration date from token: %w", err)
		}

		// Regenerate a refresh token for the user
		jwtRefreshToken, err := refreshTokenService.CreateRefreshToken(userDetails.ID)
		if err != nil {
			return fmt.Errorf("failed to create refresh token: %w", err)
		}
		if jwtRefreshToken.Equals(&entity.RefreshToken{}) {
			return fmt.Errorf("failed to create refresh token")
		}

		refreshTokenStr = jwtRefreshToken.Token

		// Update the last login time for the user
		_, err = userService.UpdateLastLogin(userDetails.ID, time.Now())
		if err != nil {
			return fmt.Errorf("failed to update last login time: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.RefreshTokenResponse{}, err
	}

	return entity.RefreshTokenResponse{
		AccessToken:    accessTokenStr,
		RefreshToken:   refreshTokenStr,
		ExpirationDate: expirationDateStr,
		TokenType:      TokenType,
	}, nil
}

// GenerateJWTToken determines the function to use for generating a JWT token based on the signing method.
// It checks the signing method from the environment variable and calls the appropriate function.
func GenerateJWTToken(user entity.User) (string, error) {
	// Load environment variables
	// LoadEnv()

	// Check the signing method from the environment variable
	if SigningMethod == jwt.SigningMethodHS256.Alg() {
		return GenerateJWTTokenWithHS256(user)
	} else if SigningMethod == jwt.SigningMethodRS256.Alg() {
		return GenerateJWTTokenWithRS256(user)
	}

	return "", fmt.Errorf("unsupported signing method: %s", SigningMethod)
}

// GenerateJWTTokenWithHS256 generates a JWT token using the HS256 signing method.
// It creates the claims for the token and signs it with the secret key from the environment variable.
func GenerateJWTTokenWithHS256(user entity.User) (string, error) {
	// Load environment variables
	// LoadEnv()

	// Set the now time
	// This is used to set the issued at (iat) and expiration (exp) claims
	now := time.Now().Unix()

	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub":      user.Username,
		"aud":      JWTAudience,
		"iss":      JWTIssuer,
		"iat":      now,
		"exp":      GetJWTExpiration(now),
		"email":    user.Email,
		"userid":   user.ID,
		"username": user.Username,
		"roles":    ExtractRoleNames(user.Roles),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// GenerateJWTTokenWithRS256 generates a JWT token using the RS256 signing method.
// It creates the claims for the token and signs it with the private key loaded from the file.
func GenerateJWTTokenWithRS256(user entity.User) (string, error) {
	// Load environment variables
	// LoadEnv()

	// Load the private key from the file
	privateKey, err := jwtutil.LoadPrivateKey()
	if err != nil {
		return "", err
	}

	// Set the now time
	// This is used to set the issued at (iat) and expiration (exp) claims
	now := time.Now().Unix()

	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub":      user.Username,
		"aud":      JWTAudience,
		"iss":      JWTIssuer,
		"iat":      now,
		"exp":      GetJWTExpiration(now),
		"email":    user.Email,
		"userid":   user.ID,
		"username": user.Username,
		"roles":    ExtractRoleNames(user.Roles),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ParseJWTToken determines the function to use for parsing a JWT token based on the signing method.
// It checks the signing method from the environment variable and calls the appropriate function.
func ParseJWTToken(tokenStr string) (*jwt.Token, error) {
	// Load environment variables
	// LoadEnv()

	// Check the signing method from the environment variable
	if SigningMethod == jwt.SigningMethodHS256.Alg() {
		return ParseJWTTokenWithHS256(tokenStr)
	} else if SigningMethod == jwt.SigningMethodRS256.Alg() {
		return ParseJWTTokenWithRS256(tokenStr)
	}

	return nil, fmt.Errorf("unsupported signing method: %s", SigningMethod)
}

// ParseJWTTokenWithHS256 parses a JWT token using the HS256 signing method.
// It validates the token and returns the parsed token object.
func ParseJWTTokenWithHS256(tokenStr string) (*jwt.Token, error) {
	// Load environment variables
	// LoadEnv()

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}
	return token, nil
}

// ParseJWTTokenWithRS256 parses a JWT token using the RS256 signing method.
// It validates the token and returns the parsed token object.
func ParseJWTTokenWithRS256(tokenStr string) (*jwt.Token, error) {
	// Load the public key from the file
	publicKey, err := jwtutil.LoadPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %v", err)
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}
	return token, nil
}

// GetRefreshTokenExpiration calculates the expiration time for the refresh token.
func GetJWTExpiration(now int64) int64 {
	// Load environment variables
	// LoadEnv()

	expHour, err := strconv.Atoi(JWTExpirationHour)
	if err != nil {
		return now + int64(time.Hour.Seconds()*24)
	}
	if expHour <= 0 {
		expHour = 24
	}

	return now + int64(time.Duration(expHour)*time.Hour/time.Second)
}

// ExtractRoleNames extracts the role names from a slice of roles.
func ExtractRoleNames(roles []entity.Role) []string {
	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = r.Name
	}
	return names
}

// GetExpirationDateFromToken extracts the expiration date from the JWT token claims.
func GetExpirationDateFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to extract claims from token")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("exp claim not found or not a float64")
	}

	expirationDate := time.Unix(int64(expFloat), 0).Format(time.RFC3339)
	return expirationDate, nil
}
