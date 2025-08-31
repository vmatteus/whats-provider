package infrastructure

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// GormMessage representa a entidade Message para GORM
type GormMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InstanceID string    `gorm:"type:varchar(255);not null;index"`
	Phone      string    `gorm:"type:varchar(20);not null"`
	Type       string    `gorm:"type:varchar(20);not null"`
	Content    string    `gorm:"type:text;not null"`
	MediaURL   *string   `gorm:"type:text"`
	Status     string    `gorm:"type:varchar(20);not null;default:'pending'"`
	ProviderID *string   `gorm:"type:varchar(255)"`
	Error      *string   `gorm:"type:text"`
	CreatedAt  int64     `gorm:"autoCreateTime"`
	UpdatedAt  int64     `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela
func (GormMessage) TableName() string {
	return "whatsapp_messages"
}

// toDomain converte GormMessage para domain.Message
func (g *GormMessage) toDomain() *domain.Message {
	return &domain.Message{
		ID:         g.ID,
		InstanceID: g.InstanceID,
		Phone:      g.Phone,
		Type:       domain.MessageType(g.Type),
		Content:    g.Content,
		MediaURL:   g.MediaURL,
		Status:     domain.MessageStatus(g.Status),
		ProviderID: g.ProviderID,
		Error:      g.Error,
		CreatedAt:  timeFromUnix(g.CreatedAt),
		UpdatedAt:  timeFromUnix(g.UpdatedAt),
	}
}

// fromDomain converte domain.Message para GormMessage
func (g *GormMessage) fromDomain(message *domain.Message) {
	g.ID = message.ID
	g.InstanceID = message.InstanceID
	g.Phone = message.Phone
	g.Type = string(message.Type)
	g.Content = message.Content
	g.MediaURL = message.MediaURL
	g.Status = string(message.Status)
	g.ProviderID = message.ProviderID
	g.Error = message.Error
	g.CreatedAt = timeToUnix(message.CreatedAt)
	g.UpdatedAt = timeToUnix(message.UpdatedAt)
}

// GormMessageRepository implementa MessageRepository usando GORM
type GormMessageRepository struct {
	db *gorm.DB
}

// NewGormMessageRepository cria um novo repositório de mensagens
func NewGormMessageRepository(db *gorm.DB) *GormMessageRepository {
	return &GormMessageRepository{db: db}
}

// Save salva uma mensagem
func (r *GormMessageRepository) Save(ctx context.Context, message *domain.Message) error {
	var gormMessage GormMessage
	gormMessage.fromDomain(message)

	if err := r.db.WithContext(ctx).Create(&gormMessage).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// GetByID obtém uma mensagem por ID
func (r *GormMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	var gormMessage GormMessage

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormMessage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return gormMessage.toDomain(), nil
}

// GetByInstanceID obtém mensagens por ID da instância
func (r *GormMessageRepository) GetByInstanceID(ctx context.Context, instanceID string, limit, offset int) ([]*domain.Message, error) {
	var gormMessages []GormMessage

	query := r.db.WithContext(ctx).Where("instance_id = ?", instanceID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&gormMessages).Error; err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	messages := make([]*domain.Message, len(gormMessages))
	for i, gormMessage := range gormMessages {
		messages[i] = gormMessage.toDomain()
	}

	return messages, nil
}

// UpdateStatus atualiza o status de uma mensagem
func (r *GormMessageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.MessageStatus, providerID *string, errorMsg *string) error {
	updates := map[string]interface{}{
		"status":     string(status),
		"updated_at": timeToUnix(timeNow()),
	}

	if providerID != nil {
		updates["provider_id"] = *providerID
	}

	if errorMsg != nil {
		updates["error"] = *errorMsg
	}

	if err := r.db.WithContext(ctx).Model(&GormMessage{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}
