package domain

import (
	"context"

	"github.com/google/uuid"
)

// WhatsAppProvider define a interface que todos os provedores devem implementar
type WhatsAppProvider interface {
	// GetName retorna o nome do provedor
	GetName() string

	// SendMessage envia uma mensagem
	SendMessage(ctx context.Context, instance *Instance, request SendMessageRequest) (*SendMessageResponse, error)

	// GetInstanceStatus verifica o status de uma instância
	GetInstanceStatus(ctx context.Context, instance *Instance) (*InstanceInfo, error)

	// CreateInstance cria uma nova instância
	CreateInstance(ctx context.Context, request CreateInstanceRequest) (*Instance, error)

	// DeleteInstance remove uma instância
	DeleteInstance(ctx context.Context, instance *Instance) error

	// ValidateToken valida se o token é válido
	ValidateToken(ctx context.Context, token string) error

	// UpdateProfileName atualiza o nome do perfil da instância
	UpdateProfileName(ctx context.Context, instance *Instance, request UpdateProfileNameRequest) (*UpdateProfileResponse, error)

	// UpdateProfilePicture atualiza a foto do perfil da instância
	UpdateProfilePicture(ctx context.Context, instance *Instance, request UpdateProfilePictureRequest) (*UpdateProfileResponse, error)
}

// MessageRepository define a interface para persistência de mensagens
type MessageRepository interface {
	Save(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*Message, error)
	GetByInstanceID(ctx context.Context, instanceID string, limit, offset int) ([]*Message, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status MessageStatus, providerID *string, errorMsg *string) error
}

// InstanceRepository define a interface para persistência de instâncias
type InstanceRepository interface {
	Save(ctx context.Context, instance *Instance) error
	GetByID(ctx context.Context, id uuid.UUID) (*Instance, error)
	GetByToken(ctx context.Context, token string) (*Instance, error)
	GetByInstanceID(ctx context.Context, instanceID string) (*Instance, error)
	GetAll(ctx context.Context) ([]*Instance, error)
	Update(ctx context.Context, instance *Instance) error
	Delete(ctx context.Context, id uuid.UUID) error
}
