package whatsapp

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/whatsapp/application"
	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
	"github.com/your-org/boilerplate-go/internal/whatsapp/infrastructure"
	"github.com/your-org/boilerplate-go/internal/whatsapp/infrastructure/providers"
	"github.com/your-org/boilerplate-go/internal/whatsapp/presentation"
)

// Module configura as dependências do módulo WhatsApp
var Module = fx.Module("whatsapp",
	// Repositórios
	fx.Provide(
		fx.Annotate(
			infrastructure.NewGormMessageRepository,
			fx.As(new(domain.MessageRepository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			infrastructure.NewGormInstanceRepository,
			fx.As(new(domain.InstanceRepository)),
		),
	),

	// Provider Factory e Registry
	fx.Provide(
		fx.Annotate(
			infrastructure.NewDefaultProviderFactory,
			fx.As(new(domain.ProviderFactory)),
		),
	),
	fx.Provide(
		fx.Annotate(
			infrastructure.NewDefaultProviderRegistry,
			fx.As(new(domain.ProviderRegistry)),
		),
	),

	// Providers individuais
	fx.Provide(newZAPIProviderWithConfig),

	// Serviços
	fx.Provide(application.NewWhatsAppService),

	// Controllers
	fx.Provide(presentation.NewWhatsAppController),

	// Configuração dos provedores
	fx.Invoke(registerProviders),
	fx.Invoke(setupProviderFactory),
)

// registerProviders registra todos os provedores no serviço
func registerProviders(
	service *application.WhatsAppService,
	zapiProvider *providers.ZAPIProvider,
) {
	err := service.RegisterProvider(zapiProvider)
	if err != nil {
		// Log error but don't panic, let the app continue
		// The logger will already log this error in the service
	}
}

// setupProviderFactory configura o factory com os criadores de providers
func setupProviderFactory(
	factory domain.ProviderFactory,
) {
	// Registra o creator do Z-API
	err := factory.RegisterProvider("z-api", providers.GetZAPIProviderCreator())
	if err != nil {
		// Log error but don't panic
	}
}

// newZAPIProviderWithConfig cria um ZAPIProvider com configuração injetada
func newZAPIProviderWithConfig(cfg *config.Config, logger zerolog.Logger) *providers.ZAPIProvider {
	zapiConfig := providers.ZAPIConfig{
		BaseURL:     cfg.WhatsApp.ZApi.BaseURL,
		ClientToken: cfg.WhatsApp.ZApi.ClientToken,
	}

	return providers.NewZAPIProviderWithConfig(zapiConfig, logger)
}
