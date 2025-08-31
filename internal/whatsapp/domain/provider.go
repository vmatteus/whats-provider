package domain

import (
	"context"
)

// ProviderConfig representa uma configuração genérica para qualquer provider
type ProviderConfig map[string]interface{}

// ProviderFeature representa uma funcionalidade suportada pelo provider
type ProviderFeature string

const (
	FeatureTextMessages   ProviderFeature = "text_messages"
	FeatureImageMessages  ProviderFeature = "image_messages"
	FeatureVideoMessages  ProviderFeature = "video_messages"
	FeatureAudioMessages  ProviderFeature = "audio_messages"
	FeatureFileMessages   ProviderFeature = "file_messages"
	FeatureGroupMessages  ProviderFeature = "group_messages"
	FeatureWebhooks       ProviderFeature = "webhooks"
	FeatureStatusCheck    ProviderFeature = "status_check"
	FeatureProfileName    ProviderFeature = "profile_name"
	FeatureProfilePicture ProviderFeature = "profile_picture"
)

// WhatsAppProviderExtended estende a interface WhatsAppProvider com funcionalidades adicionais
type WhatsAppProviderExtended interface {
	WhatsAppProvider

	// Configure configura o provider com os parâmetros específicos
	Configure(config ProviderConfig) error

	// HealthCheck verifica se o provider está funcionando
	HealthCheck(ctx context.Context) error

	// GetSupportedFeatures retorna as funcionalidades suportadas pelo provider
	GetSupportedFeatures() []ProviderFeature
}

// ProviderFactory define a interface para criação de providers
type ProviderFactory interface {
	// CreateProvider cria um novo provider com base no tipo e configuração
	CreateProvider(providerType string, config ProviderConfig) (WhatsAppProvider, error)

	// GetSupportedProviders retorna a lista de providers suportados
	GetSupportedProviders() []string

	// RegisterProvider registra um novo tipo de provider
	RegisterProvider(providerType string, creator ProviderCreator) error
}

// ProviderCreator define a função para criar um provider específico
type ProviderCreator func(config ProviderConfig) (WhatsAppProvider, error)

// ProviderRegistry gerencia o registro e criação de providers
type ProviderRegistry interface {
	// Register registra um provider no sistema
	Register(provider WhatsAppProvider) error

	// Get obtém um provider pelo nome
	Get(name string) (WhatsAppProvider, bool)

	// GetAll obtém todos os providers registrados
	GetAll() map[string]WhatsAppProvider

	// Remove remove um provider do registro
	Remove(name string) error

	// List lista os nomes de todos os providers registrados
	List() []string
}
