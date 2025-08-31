package infrastructure

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// DefaultProviderFactory implementa a ProviderFactory
type DefaultProviderFactory struct {
	creators map[string]domain.ProviderCreator
	logger   zerolog.Logger
	mu       sync.RWMutex
}

// NewDefaultProviderFactory cria uma nova instância do factory
func NewDefaultProviderFactory(logger zerolog.Logger) *DefaultProviderFactory {
	return &DefaultProviderFactory{
		creators: make(map[string]domain.ProviderCreator),
		logger:   logger.With().Str("component", "provider_factory").Logger(),
	}
}

// CreateProvider cria um novo provider com base no tipo e configuração
func (f *DefaultProviderFactory) CreateProvider(providerType string, config domain.ProviderConfig) (domain.WhatsAppProvider, error) {
	f.mu.RLock()
	creator, exists := f.creators[providerType]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider type %s not supported", providerType)
	}

	provider, err := creator(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", providerType, err)
	}

	f.logger.Info().
		Str("provider_type", providerType).
		Msg("Provider created successfully")

	return provider, nil
}

// GetSupportedProviders retorna a lista de providers suportados
func (f *DefaultProviderFactory) GetSupportedProviders() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	providers := make([]string, 0, len(f.creators))
	for providerType := range f.creators {
		providers = append(providers, providerType)
	}
	return providers
}

// RegisterProvider registra um novo tipo de provider
func (f *DefaultProviderFactory) RegisterProvider(providerType string, creator domain.ProviderCreator) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.creators[providerType]; exists {
		return fmt.Errorf("provider type %s already registered", providerType)
	}

	f.creators[providerType] = creator
	f.logger.Info().
		Str("provider_type", providerType).
		Msg("Provider type registered successfully")

	return nil
}

// DefaultProviderRegistry implementa a ProviderRegistry
type DefaultProviderRegistry struct {
	providers map[string]domain.WhatsAppProvider
	logger    zerolog.Logger
	mu        sync.RWMutex
}

// NewDefaultProviderRegistry cria uma nova instância do registry
func NewDefaultProviderRegistry(logger zerolog.Logger) *DefaultProviderRegistry {
	return &DefaultProviderRegistry{
		providers: make(map[string]domain.WhatsAppProvider),
		logger:    logger.With().Str("component", "provider_registry").Logger(),
	}
}

// Register registra um provider no sistema
func (r *DefaultProviderRegistry) Register(provider domain.WhatsAppProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.GetName()
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	r.providers[name] = provider
	r.logger.Info().
		Str("provider_name", name).
		Msg("Provider registered successfully")

	return nil
}

// Get obtém um provider pelo nome
func (r *DefaultProviderRegistry) Get(name string) (domain.WhatsAppProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	return provider, exists
}

// GetAll obtém todos os providers registrados
func (r *DefaultProviderRegistry) GetAll() map[string]domain.WhatsAppProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Cria uma cópia para evitar race conditions
	result := make(map[string]domain.WhatsAppProvider, len(r.providers))
	for name, provider := range r.providers {
		result[name] = provider
	}
	return result
}

// Remove remove um provider do registro
func (r *DefaultProviderRegistry) Remove(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}

	delete(r.providers, name)
	r.logger.Info().
		Str("provider_name", name).
		Msg("Provider removed successfully")

	return nil
}

// List lista os nomes de todos os providers registrados
func (r *DefaultProviderRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
