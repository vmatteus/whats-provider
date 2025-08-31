package infrastructure

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/boilerplate-go/internal/whatsapp/domain"
)

// GormInstance representa a entidade Instance para GORM
type GormInstance struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Phone      *string   `gorm:"type:varchar(20)"`
	Status     string    `gorm:"type:varchar(20);not null;default:'disconnected'"`
	Provider   string    `gorm:"type:varchar(50);not null"`
	InstanceID string    `gorm:"type:varchar(255);not null"`
	Token      string    `gorm:"type:varchar(255);not null"`
	Config     string    `gorm:"type:jsonb"`
	Error      *string   `gorm:"type:text"`
	CreatedAt  int64     `gorm:"autoCreateTime"`
	UpdatedAt  int64     `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela
func (GormInstance) TableName() string {
	return "whatsapp_instances"
}

// toDomain converte GormInstance para domain.Instance
func (g *GormInstance) toDomain() *domain.Instance {
	var config map[string]any
	if g.Config != "" {
		// Aqui você pode deserializar o JSON se necessário
		config = make(map[string]any)
	}

	return &domain.Instance{
		ID:         g.ID,
		Name:       g.Name,
		Phone:      g.Phone,
		Status:     domain.InstanceStatus(g.Status),
		Provider:   g.Provider,
		InstanceID: g.InstanceID,
		Token:      g.Token,
		Config:     config,
		Error:      g.Error,
		CreatedAt:  timeFromUnix(g.CreatedAt),
		UpdatedAt:  timeFromUnix(g.UpdatedAt),
	}
}

// fromDomain converte domain.Instance para GormInstance
func (g *GormInstance) fromDomain(instance *domain.Instance) {
	g.ID = instance.ID
	g.Name = instance.Name
	g.Phone = instance.Phone
	g.Status = string(instance.Status)
	g.Provider = instance.Provider
	g.InstanceID = instance.InstanceID
	g.Token = instance.Token
	g.Error = instance.Error
	g.CreatedAt = timeToUnix(instance.CreatedAt)
	g.UpdatedAt = timeToUnix(instance.UpdatedAt)

	// Serializar config para JSON se necessário
	g.Config = "{}"
} // GormInstanceRepository implementa InstanceRepository usando GORM
type GormInstanceRepository struct {
	db *gorm.DB
}

// NewGormInstanceRepository cria um novo repositório de instâncias
func NewGormInstanceRepository(db *gorm.DB) *GormInstanceRepository {
	return &GormInstanceRepository{db: db}
}

// Save salva uma instância
func (r *GormInstanceRepository) Save(ctx context.Context, instance *domain.Instance) error {
	var gormInstance GormInstance
	gormInstance.fromDomain(instance)

	if err := r.db.WithContext(ctx).Create(&gormInstance).Error; err != nil {
		return fmt.Errorf("failed to save instance: %w", err)
	}

	return nil
}

// GetByID obtém uma instância por ID
func (r *GormInstanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Instance, error) {
	var gormInstance GormInstance

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormInstance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return gormInstance.toDomain(), nil
}

// GetByToken obtém uma instância por token
func (r *GormInstanceRepository) GetByToken(ctx context.Context, token string) (*domain.Instance, error) {
	var gormInstance GormInstance

	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&gormInstance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return gormInstance.toDomain(), nil
}

// GetByInstanceID obtém uma instância por instance_id
func (r *GormInstanceRepository) GetByInstanceID(ctx context.Context, instanceID string) (*domain.Instance, error) {
	var gormInstance GormInstance

	fmt.Printf("DEBUG: GetByInstanceID called with instanceID: %s\n", instanceID)
	if err := r.db.WithContext(ctx).Where("instance_id = ?", instanceID).First(&gormInstance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return gormInstance.toDomain(), nil
}

// GetAll obtém todas as instâncias
func (r *GormInstanceRepository) GetAll(ctx context.Context) ([]*domain.Instance, error) {
	var gormInstances []GormInstance

	if err := r.db.WithContext(ctx).Find(&gormInstances).Error; err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	instances := make([]*domain.Instance, len(gormInstances))
	for i, gormInstance := range gormInstances {
		instances[i] = gormInstance.toDomain()
	}

	return instances, nil
}

// Update atualiza uma instância
func (r *GormInstanceRepository) Update(ctx context.Context, instance *domain.Instance) error {
	var gormInstance GormInstance
	gormInstance.fromDomain(instance)

	if err := r.db.WithContext(ctx).Save(&gormInstance).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	return nil
}

// Delete remove uma instância
func (r *GormInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&GormInstance{}).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	return nil
}
