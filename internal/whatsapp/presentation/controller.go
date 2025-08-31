package presentation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/your-org/boilerplate-go/internal/response"
	"github.com/your-org/boilerplate-go/internal/whatsapp/application"
	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// WhatsAppController manipula as requisições HTTP do WhatsApp
type WhatsAppController struct {
	service *application.WhatsAppService
	logger  zerolog.Logger
}

// NewWhatsAppController cria um novo controller
func NewWhatsAppController(service *application.WhatsAppService, logger zerolog.Logger) *WhatsAppController {
	return &WhatsAppController{
		service: service,
		logger:  logger.With().Str("controller", "whatsapp").Logger(),
	}
}

// GetProviders retorna os provedores disponíveis
func (c *WhatsAppController) GetProviders(ctx *gin.Context) {
	providers := c.service.GetProviders()
	response.Success(ctx, gin.H{"providers": providers})
}

// CreateInstance cria uma nova instância
func (c *WhatsAppController) CreateInstance(ctx *gin.Context) {
	var request domain.CreateInstanceRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	instance, err := c.service.CreateInstance(ctx.Request.Context(), request)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to create instance")
		response.InternalServerError(ctx, "Failed to create instance", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{Data: instance})
}

// GetInstance obtém uma instância por ID
func (c *WhatsAppController) GetInstance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid instance ID", err.Error())
		return
	}

	instance, err := c.service.GetInstance(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "Instance not found", err.Error())
		return
	}

	response.Success(ctx, instance)
}

// GetAllInstances obtém todas as instâncias
func (c *WhatsAppController) GetAllInstances(ctx *gin.Context) {
	instances, err := c.service.GetAllInstances(ctx.Request.Context())
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to get instances")
		response.InternalServerError(ctx, "Failed to get instances", err.Error())
		return
	}

	response.Success(ctx, gin.H{"instances": instances})
}

// DeleteInstance remove uma instância
func (c *WhatsAppController) DeleteInstance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid instance ID", err.Error())
		return
	}

	err = c.service.DeleteInstance(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error().Err(err).Str("instance_id", idStr).Msg("Failed to delete instance")
		response.InternalServerError(ctx, "Failed to delete instance", err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "Instance deleted successfully"})
}

// SendMessage envia uma mensagem
func (c *WhatsAppController) SendMessage(ctx *gin.Context) {
	var request domain.SendMessageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	result, err := c.service.SendMessage(ctx.Request.Context(), request)
	if err != nil {
		c.logger.Error().Err(err).Interface("request", request).Msg("Failed to send message")
		response.InternalServerError(ctx, "Failed to send message", err.Error())
		return
	}

	response.Success(ctx, result)
}

// GetMessage obtém uma mensagem por ID
func (c *WhatsAppController) GetMessage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(ctx, "Invalid message ID", err.Error())
		return
	}

	message, err := c.service.GetMessage(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "Message not found", err.Error())
		return
	}

	response.Success(ctx, message)
}

// GetMessagesByInstance obtém mensagens de uma instância
func (c *WhatsAppController) GetMessagesByInstance(ctx *gin.Context) {
	token := ctx.Param("token")

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	messages, err := c.service.GetMessagesByInstance(ctx.Request.Context(), token, limit, offset)
	if err != nil {
		c.logger.Error().Err(err).Str("token", token).Msg("Failed to get messages")
		response.InternalServerError(ctx, "Failed to get messages", err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"messages": messages,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(messages),
		},
	})
} // GetInstanceStatus verifica o status de uma instância
func (c *WhatsAppController) GetInstanceStatus(ctx *gin.Context) {
	token := ctx.Param("token")

	status, err := c.service.GetInstanceStatus(ctx.Request.Context(), token)
	if err != nil {
		c.logger.Error().Err(err).Str("token", token).Msg("Failed to get instance status")
		response.InternalServerError(ctx, "Failed to get instance status", err.Error())
		return
	}

	response.Success(ctx, status)
}

// RegisterRoutes registra as rotas do controller
func (c *WhatsAppController) RegisterRoutes(router *gin.RouterGroup) {
	whatsapp := router.Group("/whatsapp")
	{
		// Provedores
		whatsapp.GET("/providers", c.GetProviders)

		// Instâncias
		whatsapp.POST("/instances", c.CreateInstance)
		whatsapp.GET("/instances", c.GetAllInstances)
		whatsapp.GET("/instances/:id", c.GetInstance)
		whatsapp.DELETE("/instances/:id", c.DeleteInstance)

		// Status e mensagens por token (não UUID)
		whatsapp.GET("/status/:token", c.GetInstanceStatus)
		whatsapp.GET("/messages/instance/:token", c.GetMessagesByInstance)

		// Mensagens
		whatsapp.POST("/messages", c.SendMessage)
		whatsapp.GET("/messages/:id", c.GetMessage)
	}
}
