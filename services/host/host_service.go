package host

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
)

func GetHostAll() (*[]models.Host, error) {
	var hosts []models.Host
	if err := database.DB.Preload("Studio").Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("Host not found")
	}

	return &hosts, nil
}

func CreateHost(dto *dto.CreateHostDTO) error {
	host := models.Host{
		Name:     dto.Name,
		Phone:    dto.Phone,
		StudioID: dto.StudioID,
		UserID:   dto.UserID,
	}

	if err := database.DB.Create(&host).Error; err != nil {
		return fmt.Errorf("failed to create host: %w", err)
	}

	// Preload Studio setelah create
	if err := database.DB.Preload("Studio").First(&host, host.ID).Error; err != nil {
		return fmt.Errorf("failed to load studio for new host: %w", err)
	}

	return nil
}
