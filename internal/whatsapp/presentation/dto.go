package presentation

// UpdateProfileNameDTO representa o DTO para atualização de nome do perfil
type UpdateProfileNameDTO struct {
	InstanceID string `json:"instance_id" binding:"required" example:"YOUR_INSTANCE_ID"`
	Name       string `json:"name" binding:"required" example:"Vinicius"`
}

// UpdateProfilePictureDTO representa o DTO para atualização de foto do perfil
type UpdateProfilePictureDTO struct {
	InstanceID string `json:"instance_id" binding:"required" example:"YOUR_INSTANCE_ID"`
	PictureURL string `json:"picture_url" binding:"required" example:"https://app.z-api.io/logos/zapi-dark.png"`
}

// ProfileUpdateResponseDTO representa a resposta de atualização de perfil
type ProfileUpdateResponseDTO struct {
	Success bool    `json:"success" example:"true"`
	Error   *string `json:"error,omitempty" example:"null"`
}
