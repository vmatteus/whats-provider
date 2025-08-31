package domain

// UpdateProfileNameRequest representa uma requisição para atualizar o nome do perfil
type UpdateProfileNameRequest struct {
	InstanceID string `json:"instance_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

// UpdateProfilePictureRequest representa uma requisição para atualizar a foto do perfil
type UpdateProfilePictureRequest struct {
	InstanceID string `json:"instance_id" binding:"required"`
	PictureURL string `json:"picture_url" binding:"required"`
}

// UpdateProfileResponse representa a resposta de atualização de perfil
type UpdateProfileResponse struct {
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
}
