package account

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/shopee"
)

func GetAccountAll() (*[]models.Account, error) {
	var accounts []models.Account

	if err := database.DB.Preload("Studio").Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("Account not found")
	}

	return &accounts, nil
}

func GetAccountByUniqueId(uid string) (*models.Account, error) {
	var account models.Account

	if err := database.DB.Where("unique_id = ?", uid).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func CreateOrUpdateAccount(dto *dto.CreateAccountCookiesDTO) error {

	data, err := shopee.GetAccount(dto.Cookies)
	if err != nil {
		return fmt.Errorf("Expired or failed cookies")
	}

	account := models.Account{
		Name:     data["name"].(string),
		Username: data["username"].(string),
		Email:    data["email"].(string),
		UniqueID: data["unique_id"].(string),
		Platform: "Shopee",
		Cookies:  dto.Cookies,
		StudioID: dto.StudioID,
	}

	existing, err := GetAccountByUniqueId(account.UniqueID)
	if err != nil {
		if err := database.DB.Create(&account).Error; err != nil {
			return fmt.Errorf("Failed to create account")
		}
	} else {
		existing.Cookies = dto.Cookies
		if err := database.DB.Model(&existing).Updates(dto).Error; err != nil {
			return fmt.Errorf("Failed to update cookies")
		}
	}
	return nil
}
