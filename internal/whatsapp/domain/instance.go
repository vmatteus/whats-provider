package domain

import (
	"time"

	"github.com/google/uuid"
)

// InstanceStatus representa o status de uma instância
type InstanceStatus string

const (
	InstanceConnected    InstanceStatus = "connected"
	InstanceDisconnected InstanceStatus = "disconnected"
	InstanceConnecting   InstanceStatus = "connecting"
	InstanceError        InstanceStatus = "error"
)

// Instance representa uma instância do WhatsApp
type Instance struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Phone      *string        `json:"phone,omitempty"`
	Status     InstanceStatus `json:"status"`
	Provider   string         `json:"provider"`
	InstanceID string         `json:"instance_id"` // ID da instância no provedor (ex: Z-API)
	Token      string         `json:"token"`       // Token de autenticação
	Config     map[string]any `json:"config,omitempty"`
	Error      *string        `json:"error,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// CreateInstanceRequest representa uma requisição para criar instância
type CreateInstanceRequest struct {
	Name       string         `json:"name" binding:"required"`
	Provider   string         `json:"provider" binding:"required"`
	InstanceID string         `json:"instance_id" binding:"required"` // ID da instância no provedor
	Token      string         `json:"token" binding:"required"`       // Token de autenticação
	Config     map[string]any `json:"config,omitempty"`
}

// InstanceInfo representa informações da instância
type InstanceInfo struct {
	ID     uuid.UUID      `json:"id"`
	Name   string         `json:"name"`
	Phone  *string        `json:"phone,omitempty"`
	Status InstanceStatus `json:"status"`
	Error  *string        `json:"error,omitempty"`
}
