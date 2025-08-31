package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// WhatsAppService gerencia todas as operações do WhatsApp
type WhatsAppService struct {
	providers    map[string]domain.WhatsAppProvider
	messageRepo  domain.MessageRepository
	instanceRepo domain.InstanceRepository
	logger       zerolog.Logger
}

// NewWhatsAppService cria uma nova instância do serviço
func NewWhatsAppService(
	messageRepo domain.MessageRepository,
	instanceRepo domain.InstanceRepository,
	logger zerolog.Logger,
) *WhatsAppService {
	return &WhatsAppService{
		providers:    make(map[string]domain.WhatsAppProvider),
		messageRepo:  messageRepo,
		instanceRepo: instanceRepo,
		logger:       logger.With().Str("service", "whatsapp").Logger(),
	}
}

// RegisterProvider registra um novo provedor
func (s *WhatsAppService) RegisterProvider(provider domain.WhatsAppProvider) {
	s.providers[provider.GetName()] = provider
	s.logger.Info().Str("provider", provider.GetName()).Msg("Provider registered")
}

// GetProviders retorna todos os provedores registrados
func (s *WhatsAppService) GetProviders() []string {
	providers := make([]string, 0, len(s.providers))
	for name := range s.providers {
		providers = append(providers, name)
	}
	return providers
}

// CreateInstance cria uma nova instância do WhatsApp
func (s *WhatsAppService) CreateInstance(ctx context.Context, request domain.CreateInstanceRequest) (*domain.Instance, error) {
	provider, exists := s.providers[request.Provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", request.Provider)
	}

	// Valida o token
	if err := provider.ValidateToken(ctx, request.Token); err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Cria a instância no provedor
	instance, err := provider.CreateInstance(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance in provider: %w", err)
	}

	// Salva no banco de dados
	if err := s.instanceRepo.Save(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	s.logger.Info().
		Str("instance_id", instance.ID.String()).
		Str("provider", request.Provider).
		Msg("Instance created successfully")

	return instance, nil
}

// GetInstance obtém uma instância por ID
func (s *WhatsAppService) GetInstance(ctx context.Context, id uuid.UUID) (*domain.Instance, error) {
	return s.instanceRepo.GetByID(ctx, id)
}

// GetAllInstances obtém todas as instâncias
func (s *WhatsAppService) GetAllInstances(ctx context.Context) ([]*domain.Instance, error) {
	return s.instanceRepo.GetAll(ctx)
}

// DeleteInstance remove uma instância
func (s *WhatsAppService) DeleteInstance(ctx context.Context, id uuid.UUID) error {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("instance not found: %w", err)
	}

	provider, exists := s.providers[instance.Provider]
	if !exists {
		return fmt.Errorf("provider %s not found", instance.Provider)
	}

	// Remove do provedor
	if err := provider.DeleteInstance(ctx, instance); err != nil {
		s.logger.Warn().Err(err).Msg("Failed to delete instance from provider")
	}

	// Remove do banco de dados
	if err := s.instanceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	s.logger.Info().
		Str("instance_id", id.String()).
		Msg("Instance deleted successfully")

	return nil
}

// SendMessage envia uma mensagem
func (s *WhatsAppService) SendMessage(ctx context.Context, request domain.SendMessageRequest) (*domain.SendMessageResponse, error) {
	s.logger.Info().Str("instance_id", request.InstanceID).Msg("DEBUG: SendMessage called")

	// Converte o instance_id string para UUID
	instanceUUID, err := uuid.Parse(request.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("invalid instance ID: %w", err)
	}

	s.logger.Info().Str("uuid", instanceUUID.String()).Msg("DEBUG: Parsed UUID, calling GetByID")

	// Busca a instância por UUID
	instance, err := s.instanceRepo.GetByID(ctx, instanceUUID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	provider, exists := s.providers[instance.Provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", instance.Provider)
	}

	// Cria a mensagem no banco de dados
	message := &domain.Message{
		ID:         uuid.New(),
		InstanceID: request.InstanceID,
		Phone:      request.Phone,
		Type:       request.Type,
		Content:    request.Content,
		MediaURL:   request.MediaURL,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.messageRepo.Save(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	// Envia através do provedor
	response, err := provider.SendMessage(ctx, instance, request)
	if err != nil {
		// Atualiza status para erro
		errorMsg := err.Error()
		_ = s.messageRepo.UpdateStatus(ctx, message.ID, domain.StatusFailed, nil, &errorMsg)
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Atualiza o status da mensagem
	_ = s.messageRepo.UpdateStatus(ctx, message.ID, response.Status, response.ProviderID, response.Error)

	s.logger.Info().
		Str("message_id", message.ID.String()).
		Str("instance_id", request.InstanceID).
		Str("phone", request.Phone).
		Msg("Message sent successfully")

	response.ID = message.ID
	return response, nil
}

// GetMessage obtém uma mensagem por ID
func (s *WhatsAppService) GetMessage(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	return s.messageRepo.GetByID(ctx, id)
}

// GetMessagesByInstance obtém mensagens de uma instância
func (s *WhatsAppService) GetMessagesByInstance(ctx context.Context, instanceID string, limit, offset int) ([]*domain.Message, error) {
	return s.messageRepo.GetByInstanceID(ctx, instanceID, limit, offset)
}

// GetInstanceStatus verifica o status de uma instância
func (s *WhatsAppService) GetInstanceStatus(ctx context.Context, instanceID string) (*domain.InstanceInfo, error) {
	instance, err := s.instanceRepo.GetByInstanceID(ctx, instanceID)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	provider, exists := s.providers[instance.Provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", instance.Provider)
	}

	return provider.GetInstanceStatus(ctx, instance)
}
