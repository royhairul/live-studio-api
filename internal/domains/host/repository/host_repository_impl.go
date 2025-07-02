package repository

import (
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
func (h *HostRepositoryImpl) Create(host *entity.Host) (*entity.Host, error) {
	if err := h.DB.Create(host).Error; err != nil {
		return nil, err
	}

	if err := h.DB.Preload("Studio").First(host).Error; err != nil {
		return nil, err
	}

	return host, nil
}

// FindAll implements HostRepository.
func (h *HostRepositoryImpl) FindAll() ([]*entity.Host, error) {
	var hosts []*entity.Host
	if err := h.DB.Preload("Studio").Find(&hosts).Error; err != nil {
		return nil, err
	}

	return hosts, nil
}

// FindByID implements HostRepository.
func (h *HostRepositoryImpl) FindByID(id string) (*entity.Host, error) {
	var host entity.Host
	if err := h.DB.Preload("Studio").Where("id = ?", id).First(&host).Error; err != nil {
		return nil, err
	}

	return &host, nil
}

// Update implements HostRepository.
func (h *HostRepositoryImpl) Update(host *entity.Host) (*entity.Host, error) {
	if err := h.DB.Model(&entity.Host{}).Where("id = ?", host.ID).Updates(&host).Error; err != nil {
		return nil, err
	}

	if err := h.DB.Preload("Studio").First(host).Error; err != nil {
		return nil, err
	}

	return host, nil
}

// Delete implements HostRepository.
func (h *HostRepositoryImpl) Delete(id string) error {
	if err := h.DB.Delete(&entity.Host{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
