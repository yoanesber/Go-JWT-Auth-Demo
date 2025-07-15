package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/service"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
)

// This struct defines the AuthHandler which handles HTTP requests related to authentication.
// It contains a service field of type AuthService which is used to interact with the authentication data layer.
type AuthHandler struct {
	Service service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
// It initializes the AuthHandler struct with the provided AuthService.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{Service: authService}
}

// Login handles user login requests.
// It validates the request, authenticates the user, and returns a JWT token if successful.
// @Summary      User login
// @Description  User login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      Auth  true  "Login request"
// @Success      200  {object}  model.HttpResponse for successful login
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      401  {object}  model.HttpResponse for unauthorized
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Bind the request body to the LoginRequest struct
	// This struct contains the username and password fields
	var loginReq entity.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		httputil.BadRequest(c, "Invalid request", err.Error())
		return
	}

	// Call the service to authenticate the user and get the token
	loginResp, err := h.Service.Login(loginReq)

	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			httputil.BadRequestMap(c, "Failed to login", validation.FormatValidationErrors(err))
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.Unauthorized(c, "Invalid credentials", "Username or password is incorrect")
			return
		}

		httputil.Unauthorized(c, "Failed to login", err.Error())
		return
	}

	httputil.Success(c, "Login successful", loginResp)
}

// RefreshToken handles token refresh requests.
// It validates the request, checks the refresh token, and returns a new JWT token if successful.
// @Summary      Refresh token
// @Description  Refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      entity.RefreshTokenRequest  true  "Refresh token request"
// @Success      200  {object}  model.HttpResponse for successful token refresh
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      401  {object}  model.HttpResponse for unauthorized
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Bind the request body to the RefreshTokenRequest struct
	// This struct contains the refresh token field
	var refreshTokenReq entity.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenReq); err != nil {
		httputil.BadRequest(c, "Invalid request", err.Error())
		return
	}

	// Call the service to refresh the token
	refreshTokenResp, err := h.Service.RefreshToken(refreshTokenReq)

	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			httputil.BadRequestMap(c, "Failed to refresh token", validation.FormatValidationErrors(err))
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.Unauthorized(c, "Invalid refresh token", "Refresh token is invalid")
			return
		}

		// Handle other errors, such as database connection issues
		// or query execution errors
		httputil.Unauthorized(c, "Failed to refresh token", err.Error())
		return
	}

	httputil.Success(c, "Token refreshed successfully", refreshTokenResp)
}
