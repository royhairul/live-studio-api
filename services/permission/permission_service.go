package permission

import (
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/models"
)

func GetAllPermission() (*[]models.Permission, error) {
	var permissions []models.Permission

	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return &permissions, nil
}

func GetAllPermissionGrouped() (*map[string][]models.Permission, error) {
	var permissions []models.Permission

	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Permission)

	for _, perm := range permissions {
		grouped[perm.Group] = append(grouped[perm.Group], perm)
	}

	return &grouped, nil
}
