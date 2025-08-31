package domain

import (
	"time"

	"github.com/google/uuid"
)

// MessageType representa os tipos de mensagem suportados
type MessageType string

const (
	TextMessage     MessageType = "text"
	ImageMessage    MessageType = "image"
	DocumentMessage MessageType = "document"
	AudioMessage    MessageType = "audio"
	VideoMessage    MessageType = "video"
)

// MessageStatus representa o status de uma mensagem
type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// Message representa uma mensagem do WhatsApp
type Message struct {
	ID         uuid.UUID     `json:"id"`
	InstanceID string        `json:"instance_id"`
	Phone      string        `json:"phone"`
	Type       MessageType   `json:"type"`
	Content    string        `json:"content"`
	MediaURL   *string       `json:"media_url,omitempty"`
	Status     MessageStatus `json:"status"`
	ProviderID *string       `json:"provider_id,omitempty"`
	Error      *string       `json:"error,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// SendMessageRequest representa uma requisição para enviar mensagem
type SendMessageRequest struct {
	InstanceID string      `json:"instance_id" binding:"required"`
	Phone      string      `json:"phone" binding:"required"`
	Type       MessageType `json:"type" binding:"required"`
	Content    string      `json:"content" binding:"required"`
	MediaURL   *string     `json:"media_url,omitempty"`
}

// SendMessageResponse representa a resposta de envio de mensagem
type SendMessageResponse struct {
	ID         uuid.UUID     `json:"id"`
	Status     MessageStatus `json:"status"`
	ProviderID *string       `json:"provider_id,omitempty"`
	Error      *string       `json:"error,omitempty"`
}
