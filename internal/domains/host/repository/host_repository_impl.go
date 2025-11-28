package repository

import (
	"context"

	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
	"gorm.io/gorm"
)

type HostRepositoryImpl struct {
	DB *gorm.DB
}

func NewHostRepository(db *gorm.DB) HostRepository {
	return &HostRepositoryImpl{DB: db}
}

// Create implements HostRepository.
func (h *HostRepositoryImpl) Create(ctx context.Context, host *entity.Host) (*entity.Host, error) {
	if err := h.DB.WithContext(ctx).Create(host).Error; err != nil {
		return nil, err
	}

	if err := h.DB.WithContext(ctx).Preload("Studio").First(host).Error; err != nil {
		return nil, err
	}

	return host, nil
}

// FindAll implements HostRepository.
func (h *HostRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Host, error) {
	var hosts []*entity.Host
	if err := h.DB.WithContext(ctx).Preload("Studio").Find(&hosts).Error; err != nil {
		return nil, err
	}

	return hosts, nil
}

// FindByID implements HostRepository.
func (h *HostRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Host, error) {
	var host entity.Host
	if err := h.DB.WithContext(ctx).Preload("Studio").Where("id = ?", id).First(&host).Error; err != nil {
		return nil, err
	}

	return &host, nil
}

// Update implements HostRepository.
func (h *HostRepositoryImpl) Update(ctx context.Context, host *entity.Host) (*entity.Host, error) {
	if err := h.DB.WithContext(ctx).Model(&entity.Host{}).Where("id = ?", host.ID).Updates(&host).Error; err != nil {
		return nil, err
	}

	if err := h.DB.WithContext(ctx).Preload("Studio").First(host).Error; err != nil {
		return nil, err
	}

	return host, nil
}

// Delete implements HostRepository.
func (h *HostRepositoryImpl) Delete(ctx context.Context, id string) error {
	if err := h.DB.WithContext(ctx).Delete(&entity.Host{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
