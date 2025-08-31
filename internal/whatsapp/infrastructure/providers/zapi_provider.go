package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// ZAPIProvider implementa a interface WhatsAppProvider para a Z-API
type ZAPIProvider struct {
	baseURL     string
	clientToken string
	httpClient  *http.Client
	logger      zerolog.Logger
}

// ZAPIConfig representa a configuração da Z-API
type ZAPIConfig struct {
	BaseURL     string `json:"base_url"`
	ClientToken string `json:"client_token"`
}

// ZAPISendMessageRequest representa a requisição de envio da Z-API
type ZAPISendMessageRequest struct {
	Phone        string `json:"phone"`
	Message      string `json:"message,omitempty"`      // Para mensagens de texto
	DelayMessage int    `json:"delayMessage,omitempty"` // Delay para envio
	Image        string `json:"image,omitempty"`
	Video        string `json:"video,omitempty"`
	Audio        string `json:"audio,omitempty"`
}

// ZAPISendMessageResponse representa a resposta de envio da Z-API
type ZAPISendMessageResponse struct {
	ZaapID    string `json:"zaapId,omitempty"`
	MessageID string `json:"messageId,omitempty"`
	ID        string `json:"id,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ZAPIInstanceStatusResponse representa a resposta de status da instância
type ZAPIInstanceStatusResponse struct {
	Connected bool   `json:"connected"`
	Phone     string `json:"phone,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ZAPIStatusResponse representa a resposta de status da Z-API
type ZAPIStatusResponse struct {
	Status string `json:"status"`
	Phone  string `json:"phone,omitempty"`
	Error  string `json:"error,omitempty"`
}

// ZAPIUpdateProfileNameRequest representa a requisição para atualizar nome do perfil na Z-API
type ZAPIUpdateProfileNameRequest struct {
	Value string `json:"value"`
}

// ZAPIUpdateProfilePictureRequest representa a requisição para atualizar foto do perfil na Z-API
type ZAPIUpdateProfilePictureRequest struct {
	Value string `json:"value"`
}

// ZAPIUpdateProfileResponse representa a resposta de atualização de perfil da Z-API
type ZAPIUpdateProfileResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewZAPIProvider cria um novo provedor Z-API
func NewZAPIProvider(logger zerolog.Logger) *ZAPIProvider {
	return &ZAPIProvider{
		baseURL:     "https://api.z-api.io/instances",
		clientToken: "123", // TODO: vir da configuração
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.With().Str("provider", "z-api").Logger(),
	}
}

// NewZAPIProviderWithConfig cria um novo provedor Z-API com configuração personalizada
func NewZAPIProviderWithConfig(config ZAPIConfig, logger zerolog.Logger) *ZAPIProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.z-api.io/instances"
	}

	clientToken := config.ClientToken
	if clientToken == "" {
		clientToken = "123" // valor padrão
	}

	return &ZAPIProvider{
		baseURL:     baseURL,
		clientToken: clientToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.With().Str("provider", "z-api").Logger(),
	}
}

// GetName retorna o nome do provedor
func (z *ZAPIProvider) GetName() string {
	return "z-api"
}

// SendMessage envia uma mensagem através da Z-API
func (z *ZAPIProvider) SendMessage(ctx context.Context, instance *domain.Instance, request domain.SendMessageRequest) (*domain.SendMessageResponse, error) {
	// Prepara a requisição para Z-API conforme documentação
	var zapiRequest ZAPISendMessageRequest
	zapiRequest.Phone = request.Phone

	// Configura o tipo de mensagem - para texto simples
	switch request.Type {
	case domain.TextMessage:
		zapiRequest.Message = request.Content
		zapiRequest.DelayMessage = 15 // Delay padrão de 15 segundos
	case domain.ImageMessage:
		zapiRequest.Image = *request.MediaURL
		zapiRequest.Message = request.Content // Legenda
	case domain.VideoMessage:
		zapiRequest.Video = *request.MediaURL
		zapiRequest.Message = request.Content // Legenda
	case domain.AudioMessage:
		zapiRequest.Audio = *request.MediaURL
	default:
		return nil, fmt.Errorf("message type %s not supported by Z-API", request.Type)
	}

	// Monta a URL no formato correto: instances/{instance_id}/token/{token}/send-text
	var endpoint string
	switch request.Type {
	case domain.TextMessage:
		endpoint = "send-text"
	case domain.ImageMessage:
		endpoint = "send-image"
	case domain.VideoMessage:
		endpoint = "send-video"
	case domain.AudioMessage:
		endpoint = "send-audio"
	case domain.DocumentMessage:
		endpoint = "send-document"
	}

	url := fmt.Sprintf("%s/%s/token/%s/%s", z.baseURL, instance.InstanceID, instance.Token, endpoint)

	// Faz a requisição
	response, err := z.makeRequest(ctx, "POST", url, zapiRequest)
	if err != nil {
		return nil, err
	}

	var zapiResponse ZAPISendMessageResponse
	if err := json.Unmarshal(response, &zapiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Z-API response: %w", err)
	}

	// Se não há ID da mensagem, algo deu errado
	if zapiResponse.ID == "" && zapiResponse.MessageID == "" {
		errorMsg := zapiResponse.Error
		if errorMsg == "" {
			errorMsg = "unknown error from Z-API"
		}
		return &domain.SendMessageResponse{
			Status: domain.StatusFailed,
			Error:  &errorMsg,
		}, nil
	}

	// Usa o ID preferencial (messageId ou id)
	messageID := zapiResponse.MessageID
	if messageID == "" {
		messageID = zapiResponse.ID
	}

	return &domain.SendMessageResponse{
		Status:     domain.StatusSent,
		ProviderID: &messageID,
	}, nil
}

// GetInstanceStatus obtém o status de uma instância Z-API
func (z *ZAPIProvider) GetInstanceStatus(ctx context.Context, instance *domain.Instance) (*domain.InstanceInfo, error) {
	url := fmt.Sprintf("%s/%s/token/%s/status", z.baseURL, instance.InstanceID, instance.Token)

	z.logger.Info().
		Str("url", url).
		Str("instance_id", instance.InstanceID).
		Str("token", instance.Token).
		Msg("Getting instance status from Z-API")

	response, err := z.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		z.logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to get instance status from Z-API")
		return nil, err
	}

	z.logger.Debug().
		Str("response", string(response)).
		Msg("Z-API status response")

	var zapiResponse ZAPIStatusResponse
	if err := json.Unmarshal(response, &zapiResponse); err != nil {
		z.logger.Error().
			Err(err).
			Str("response", string(response)).
			Msg("Failed to parse Z-API status response")
		return nil, fmt.Errorf("failed to parse Z-API response: %w", err)
	}

	z.logger.Info().
		Str("status", zapiResponse.Status).
		Str("phone", zapiResponse.Phone).
		Str("error", zapiResponse.Error).
		Msg("Parsed Z-API status response")

	// Mapeia o status da Z-API para o nosso domínio
	var status domain.InstanceStatus
	switch zapiResponse.Status {
	case "open":
		status = domain.InstanceConnected
	case "closed":
		status = domain.InstanceDisconnected
	case "connecting":
		status = domain.InstanceConnecting
	default:
		z.logger.Warn().
			Str("received_status", zapiResponse.Status).
			Msg("Unknown status from Z-API, setting as error")
		status = domain.InstanceError
	}

	// Converte phone string para *string
	var phone *string
	if zapiResponse.Phone != "" {
		phone = &zapiResponse.Phone
	}

	// Se há erro na resposta da Z-API, definir como erro
	var errorMsg *string
	if zapiResponse.Error != "" {
		// Alguns "erros" da Z-API são na verdade estados válidos
		switch zapiResponse.Error {
		case "You are already connected.":
			// Instância já conectada - isso é um estado válido, não um erro
			z.logger.Info().
				Str("message", zapiResponse.Error).
				Msg("Z-API instance is already connected")
			status = domain.InstanceConnected
		default:
			// Outros erros são realmente problemas
			errorMsg = &zapiResponse.Error
			status = domain.InstanceError
			z.logger.Warn().
				Str("error", zapiResponse.Error).
				Msg("Z-API returned error")
		}
	}

	return &domain.InstanceInfo{
		ID:     uuid.New(), // TODO: buscar ID real da instância
		Name:   "Z-API Instance",
		Phone:  phone,
		Status: status,
		Error:  errorMsg,
	}, nil
}

// CreateInstance cria uma nova instância Z-API
func (z *ZAPIProvider) CreateInstance(ctx context.Context, request domain.CreateInstanceRequest) (*domain.Instance, error) {
	instance := &domain.Instance{
		ID:         uuid.New(),
		Name:       request.Name,
		Provider:   z.GetName(),
		InstanceID: request.InstanceID,
		Token:      request.Token,
		Config:     request.Config,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Verifica o status da instância na Z-API
	instanceInfo, err := z.GetInstanceStatus(ctx, instance)
	if err != nil {
		instance.Status = domain.InstanceError
		errorMsg := err.Error()
		instance.Error = &errorMsg
	} else {
		instance.Status = instanceInfo.Status
		instance.Phone = instanceInfo.Phone
		instance.Error = instanceInfo.Error
	}

	return instance, nil
}

// DeleteInstance remove uma instância Z-API
func (z *ZAPIProvider) DeleteInstance(ctx context.Context, instance *domain.Instance) error {
	// Z-API não tem endpoint específico para deletar instância
	// Apenas validamos se a instância existe
	_, err := z.GetInstanceStatus(ctx, instance)
	return err
}

// ValidateToken valida se o token é válido
func (z *ZAPIProvider) ValidateToken(ctx context.Context, token string) error {
	// Para validar um token, precisaríamos da instância completa
	// Por enquanto, vamos retornar nil (sempre válido)
	// TODO: Implementar validação adequada quando necessário
	return nil
}

// makeRequest faz uma requisição HTTP para a Z-API
func (z *ZAPIProvider) makeRequest(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Client-Token", z.clientToken)

	z.logger.Debug().
		Str("method", method).
		Str("url", url).
		Str("client_token", z.clientToken).
		Msg("Making request to Z-API")

	resp, err := z.httpClient.Do(req)
	if err != nil {
		z.logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to make HTTP request to Z-API")
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	z.logger.Debug().
		Int("status_code", resp.StatusCode).
		Str("response_body", string(responseBody)).
		Str("url", url).
		Msg("Z-API response received")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		z.logger.Error().
			Int("status_code", resp.StatusCode).
			Str("response_body", string(responseBody)).
			Str("url", url).
			Msg("Z-API returned error status")
		return nil, fmt.Errorf("Z-API returned error status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

// Configure configura o provider com os parâmetros específicos
func (z *ZAPIProvider) Configure(config domain.ProviderConfig) error {
	if baseURL, ok := config["base_url"].(string); ok && baseURL != "" {
		z.baseURL = baseURL
	}

	if clientToken, ok := config["client_token"].(string); ok && clientToken != "" {
		z.clientToken = clientToken
	}

	if timeout, ok := config["timeout"].(time.Duration); ok && timeout > 0 {
		z.httpClient.Timeout = timeout
	}

	z.logger.Info().Msg("Z-API provider configured successfully")
	return nil
}

// HealthCheck verifica se o provider está funcionando
func (z *ZAPIProvider) HealthCheck(ctx context.Context) error {
	// Faz uma requisição simples para verificar se a API está respondendo
	url := fmt.Sprintf("%s/status", z.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("Z-API server error: status %d", resp.StatusCode)
	}

	return nil
}

// GetSupportedFeatures retorna as funcionalidades suportadas pelo provider
func (z *ZAPIProvider) GetSupportedFeatures() []domain.ProviderFeature {
	return []domain.ProviderFeature{
		domain.FeatureTextMessages,
		domain.FeatureImageMessages,
		domain.FeatureVideoMessages,
		domain.FeatureAudioMessages,
		domain.FeatureFileMessages,
		domain.FeatureStatusCheck,
		domain.FeatureWebhooks,
		domain.FeatureProfileName,
		domain.FeatureProfilePicture,
	}
}

// UpdateProfileName atualiza o nome do perfil da instância
func (z *ZAPIProvider) UpdateProfileName(ctx context.Context, instance *domain.Instance, request domain.UpdateProfileNameRequest) (*domain.UpdateProfileResponse, error) {
	zapiRequest := ZAPIUpdateProfileNameRequest{
		Value: request.Name,
	}

	url := fmt.Sprintf("%s/%s/token/%s/profile-name", z.baseURL, instance.InstanceID, instance.Token)

	z.logger.Info().
		Str("url", url).
		Str("instance_id", instance.InstanceID).
		Str("name", request.Name).
		Msg("Updating profile name via Z-API")

	response, err := z.makeRequest(ctx, "PUT", url, zapiRequest)
	if err != nil {
		z.logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to update profile name via Z-API")

		errorMsg := err.Error()
		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &errorMsg,
		}, nil
	}

	var zapiResponse ZAPIUpdateProfileResponse
	if err := json.Unmarshal(response, &zapiResponse); err != nil {
		z.logger.Error().
			Err(err).
			Str("response", string(response)).
			Msg("Failed to parse Z-API profile name update response")

		errorMsg := fmt.Sprintf("failed to parse response: %v", err)
		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &errorMsg,
		}, nil
	}

	if zapiResponse.Error != "" {
		z.logger.Error().
			Str("error", zapiResponse.Error).
			Msg("Z-API returned error when updating profile name")

		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &zapiResponse.Error,
		}, nil
	}

	z.logger.Info().
		Str("instance_id", instance.InstanceID).
		Str("name", request.Name).
		Msg("Profile name updated successfully")

	return &domain.UpdateProfileResponse{
		Success: true,
	}, nil
}

// UpdateProfilePicture atualiza a foto do perfil da instância
func (z *ZAPIProvider) UpdateProfilePicture(ctx context.Context, instance *domain.Instance, request domain.UpdateProfilePictureRequest) (*domain.UpdateProfileResponse, error) {
	zapiRequest := ZAPIUpdateProfilePictureRequest{
		Value: request.PictureURL,
	}

	url := fmt.Sprintf("%s/%s/token/%s/profile-picture", z.baseURL, instance.InstanceID, instance.Token)

	z.logger.Info().
		Str("url", url).
		Str("instance_id", instance.InstanceID).
		Str("picture_url", request.PictureURL).
		Msg("Updating profile picture via Z-API")

	response, err := z.makeRequest(ctx, "PUT", url, zapiRequest)
	if err != nil {
		z.logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to update profile picture via Z-API")

		errorMsg := err.Error()
		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &errorMsg,
		}, nil
	}

	var zapiResponse ZAPIUpdateProfileResponse
	if err := json.Unmarshal(response, &zapiResponse); err != nil {
		z.logger.Error().
			Err(err).
			Str("response", string(response)).
			Msg("Failed to parse Z-API profile picture update response")

		errorMsg := fmt.Sprintf("failed to parse response: %v", err)
		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &errorMsg,
		}, nil
	}

	if zapiResponse.Error != "" {
		z.logger.Error().
			Str("error", zapiResponse.Error).
			Msg("Z-API returned error when updating profile picture")

		return &domain.UpdateProfileResponse{
			Success: false,
			Error:   &zapiResponse.Error,
		}, nil
	}

	z.logger.Info().
		Str("instance_id", instance.InstanceID).
		Str("picture_url", request.PictureURL).
		Msg("Profile picture updated successfully")

	return &domain.UpdateProfileResponse{
		Success: true,
	}, nil
}
