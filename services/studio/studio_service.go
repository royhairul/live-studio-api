package studio

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
)

func GetStudioAll() ([]models.Studio, error) {
	var studios []models.Studio
	if err := database.DB.Find(&studios).Error; err != nil {
		return nil, err
	}

	return studios, nil
}

func GetStudioByID(id string) (*models.Studio, error) {
	var studio models.Studio
	if err := database.DB.Where("id = ?", id).First(&studio).Error; err != nil {
		return nil, fmt.Errorf("Studio not found")
	}

	return &studio, nil
}

func UpdateStudio(id string, dto *dto.Studio) error {
	studio, err := GetStudioByID(id)
	if err != nil {
		return err
	}

	studio.Name = dto.Name
	studio.Address = dto.Address

	if err := database.DB.Updates(&studio).Error; err != nil {
		return fmt.Errorf("Failed to update studio")
	}

	return nil
}
