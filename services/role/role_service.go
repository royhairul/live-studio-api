package role

import (
	"errors"
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"gorm.io/gorm"
)

func GetAllRole() (*[]models.Role, error) {
	var roles []models.Role

	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	return &roles, nil
}

func GetRoleByID(id string) (*models.Role, error) {
	var role models.Role

	if err := database.DB.Preload("Permissions").Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func GetRoleByName(name string) (*models.Role, error) {
	var role models.Role

	if err := database.DB.Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		// Cek jika record tidak ditemukan, return nil dan error yang jelas
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role '%s' not found", name)
		}
		return nil, err
	}

	return &role, nil
}

func CreateRole(dto *dto.CreateRoleDTO) error {
	role := models.Role{
		Name: dto.Name,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		return err
	}

	if len(dto.Permissions) > 0 {
		var permissions []models.Permission
		if err := database.DB.Where("id IN ?", dto.Permissions).Find(&permissions).Error; err != nil {
			return err
		}

		if err := database.DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
			return err
		}
	}

	return nil
}

func UpdateRole(id string, dto *dto.UpdateRoleDTO) error {
	role, err := GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("Failed role not found")
	}

	role.Name = *dto.Name

	if err := database.DB.Updates(&role).Error; err != nil {
		return err
	}

	// Update relasi Permissions (replace semuanya)
	var permissions []models.Permission
	if dto.Permissions != nil && len(*dto.Permissions) > 0 {
		if err := database.DB.Where("id IN ?", *dto.Permissions).Find(&permissions).Error; err != nil {
			return err
		}
	}

	// Replace relasi many2many di tabel pivot
	if err := database.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		return err
	}

	return nil
}

func DeleteRole(id string) error {
	role, err := GetRoleByID(id)
	if err != nil {
		return fmt.Errorf("Failed role not found")
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}
