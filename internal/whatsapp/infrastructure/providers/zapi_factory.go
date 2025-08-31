package providers

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// CreateZAPIProvider cria um provider Z-API através da factory
func CreateZAPIProvider(config domain.ProviderConfig) (domain.WhatsAppProvider, error) {
	zapiConfig := ZAPIConfig{}

	// Extrai configurações do ProviderConfig
	if baseURL, ok := config["base_url"].(string); ok {
		zapiConfig.BaseURL = baseURL
	}

	if clientToken, ok := config["client_token"].(string); ok {
		zapiConfig.ClientToken = clientToken
	}

	// Se não tiver logger no config, cria um padrão
	var logger zerolog.Logger
	if loggerInterface, ok := config["logger"]; ok {
		if zlogger, ok := loggerInterface.(zerolog.Logger); ok {
			logger = zlogger
		} else {
			logger = zerolog.New(nil).With().Str("provider", "z-api").Logger()
		}
	} else {
		logger = zerolog.New(nil).With().Str("provider", "z-api").Logger()
	}

	provider := NewZAPIProviderWithConfig(zapiConfig, logger)

	// Aplica configurações adicionais
	if err := provider.Configure(config); err != nil {
		return nil, fmt.Errorf("failed to configure Z-API provider: %w", err)
	}

	return provider, nil
}

// GetZAPIProviderCreator retorna a função creator para o Z-API provider
func GetZAPIProviderCreator() domain.ProviderCreator {
	return CreateZAPIProvider
}
