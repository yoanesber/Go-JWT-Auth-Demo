package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"

	"github.com/yoanesber/go-consumer-api-with-jwt/internal/entity"
	"github.com/yoanesber/go-consumer-api-with-jwt/internal/service"
	httputil "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/http-util"
	validation "github.com/yoanesber/go-consumer-api-with-jwt/pkg/util/validation-util"
)

// This struct defines the ConsumerHandler which handles HTTP requests related to consumers.
// It contains a service field of type ConsumerService which is used to interact with the consumer data layer.
type ConsumerHandler struct {
	Service service.ConsumerService
}

// NewConsumerHandler creates a new instance of ConsumerHandler.
// It initializes the ConsumerHandler struct with the provided ConsumerService.
func NewConsumerHandler(consumerService service.ConsumerService) *ConsumerHandler {
	return &ConsumerHandler{Service: consumerService}
}

// GetAllConsumers retrieves all consumers from the database and returns them as JSON.
// @Summary      Get all consumers
// @Description  Get all consumers from the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        page   query     string  false "Page number (default is 1)"
// @Param        limit  query     string  false "Number of transactions per page (default is 10)"
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers [get]
func (h *ConsumerHandler) GetAllConsumers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httputil.BadRequest(c, "Invalid page number", "Page must be a positive integer")
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httputil.BadRequest(c, "Invalid limit", "Limit must be a positive integer")
		return
	}

	consumers, err := h.Service.GetAllConsumers(page, limit)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve consumers", err.Error())
		return
	}

	if len(consumers) == 0 {
		httputil.NotFound(c, "No consumers found", "No consumers available in the database")
		return
	}

	httputil.Success(c, "All consumers retrieved successfully", consumers)
}

// GetConsumerByID retrieves a consumer by its ID from the database and returns it as JSON.
// @Summary      Get consumer by ID
// @Description  Get a consumer by its ID from the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Consumer ID"
// @Success      200  {object}  model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers/{id} [get]
func (h *ConsumerHandler) GetConsumerByID(c *gin.Context) {
	// Parse the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		httputil.BadRequest(c, "Invalid ID", "ID cannot be empty")
		return
	}

	// Retrieve the consumer by ID from the service
	consumer, err := h.Service.GetConsumerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.NotFound(c, "Consumer not found", "No consumer found with the given ID")
			return
		}

		// If the error is not a record not found error, return a generic internal server error
		// This is to avoid exposing internal details of the error
		httputil.InternalServerError(c, "Failed to retrieve consumer", err.Error())
		return
	}

	httputil.Success(c, "Consumer retrieved successfully", consumer)
}

// GetActiveConsumers retrieves all active consumers from the database and returns them as JSON.
// @Summary      Get active consumers
// @Description  Get all active consumers from the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        page   query     string  false "Page number (default is 1)"
// @Param        limit  query     string  false "Number of transactions per page (default is 10)"
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers/active [get]
func (h *ConsumerHandler) GetActiveConsumers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httputil.BadRequest(c, "Invalid page number", "Page must be a positive integer")
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httputil.BadRequest(c, "Invalid limit", "Limit must be a positive integer")
		return
	}

	activeConsumers, err := h.Service.GetActiveConsumers(page, limit)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve active consumers", err.Error())
		return
	}

	if len(activeConsumers) == 0 {
		httputil.NotFound(c, "No active consumers found", "No active consumers available in the database")
		return
	}

	httputil.Success(c, "Active consumers retrieved successfully", activeConsumers)
}

// GetInactiveConsumers retrieves all inactive consumers from the database and returns them as JSON.
// @Summary      Get inactive consumers
// @Description  Get all inactive consumers from the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        page   query     string  false "Page number (default is 1)"
// @Param        limit  query     string  false "Number of transactions per page (default is 10)"
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers/inactive [get]
func (h *ConsumerHandler) GetInactiveConsumers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httputil.BadRequest(c, "Invalid page number", "Page must be a positive integer")
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httputil.BadRequest(c, "Invalid limit", "Limit must be a positive integer")
		return
	}

	inactiveConsumers, err := h.Service.GetInactiveConsumers(page, limit)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve inactive consumers", err.Error())
		return
	}

	if len(inactiveConsumers) == 0 {
		httputil.NotFound(c, "No inactive consumers found", "No inactive consumers available in the database")
		return
	}

	httputil.Success(c, "Inactive consumers retrieved successfully", inactiveConsumers)
}

// GetSuspendedConsumers retrieves all suspended consumers from the database and returns them as JSON.
// @Summary      Get suspended consumers
// @Description  Get all suspended consumers from the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        page   query     string  false "Page number (default is 1)"
// @Param        limit  query     string  false "Number of transactions per page (default is 10)"
// @Success      200  {array}   model.HttpResponse for successful retrieval
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers/suspended [get]
func (h *ConsumerHandler) GetSuspendedConsumers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httputil.BadRequest(c, "Invalid page number", "Page must be a positive integer")
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httputil.BadRequest(c, "Invalid limit", "Limit must be a positive integer")
		return
	}

	suspendedConsumers, err := h.Service.GetSuspendedConsumers(page, limit)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve suspended consumers", err.Error())
		return
	}

	if len(suspendedConsumers) == 0 {
		httputil.NotFound(c, "No suspended consumers found", "No suspended consumers available in the database")
		return
	}

	httputil.Success(c, "Suspended consumers retrieved successfully", suspendedConsumers)
}

// CreateConsumer creates a new consumer in the database and returns it as JSON.
// @Summary      Create consumer
// @Description  Create a new consumer in the database
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        consumer  body      Consumer  true  "Consumer object"
// @Success      201  {object}  model.HttpResponse for successful creation
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers [post]
func (h *ConsumerHandler) CreateConsumer(c *gin.Context) {
	// Bind the JSON request body to the Consumer struct
	// This will automatically validate the request body against the struct tags
	var consumer entity.Consumer
	if err := c.ShouldBindJSON(&consumer); err != nil {
		httputil.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Create the consumer using the service
	createdConsumer, err := h.Service.CreateConsumer(consumer)
	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			httputil.BadRequestMap(c, "Failed to create consumer", validation.FormatValidationErrors(err))
			return
		}

		// If the error is not a validation error, return a generic internal server error
		// This is to avoid exposing internal details of the error
		httputil.InternalServerError(c, "Failed to create consumer", err.Error())
		return
	}

	httputil.Created(c, "Consumer created successfully", createdConsumer)
}

// UpdateConsumerStatus updates the status of a consumer by its ID and returns the updated consumer as JSON.
// @Summary      Update consumer status
// @Description  Update the status of a consumer by its ID
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Consumer ID"
// @Param        status query     string  true  "New status (active, inactive, suspended)"
// @Success      200  {object}  model.HttpResponse for successful update
// @Failure      400  {object}  model.HttpResponse for bad request
// @Failure      404  {object}  model.HttpResponse for not found
// @Failure      500  {object}  model.HttpResponse for internal server error
// @Router       /consumers/{id}?status={status} [patch]
func (h *ConsumerHandler) UpdateConsumerStatus(c *gin.Context) {
	// Get the ID and status from the URL parameters
	id := c.Param("id")
	status := c.DefaultQuery("status", "")

	// Validate the ID
	if id == "" {
		httputil.BadRequest(c, "Invalid ID", "ID cannot be empty")
		return
	}

	// Validate the status
	if status != entity.ConsumerStatusActive && status != entity.ConsumerStatusInactive && status != entity.ConsumerStatusSuspended {
		httputil.BadRequest(c, "Invalid status", "Status must be one of: active, inactive, suspended")
		return
	}

	// Update the consumer status using the service
	updatedConsumer, err := h.Service.UpdateConsumerStatus(id, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httputil.NotFound(c, "Consumer not found", "No consumer found with the given ID")
			return
		}

		// If the error is not a record not found error, return a generic internal server error
		// This is to avoid exposing internal details of the error
		httputil.InternalServerError(c, "Failed to update consumer status", err.Error())
		return
	}

	httputil.Success(c, "Consumer status updated successfully", updatedConsumer)
}
