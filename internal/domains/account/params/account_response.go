package params

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/account/entity"
)

type AccountResponse struct {
	ID         uint   `json:"id"`
	UniqueID   string `json:"unique_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Platform   string `json:"platform"`
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	Device     string `json:"device"`
	Cookie     string `json:"cookie"`
	IsActive   bool   `json:"is_active"`
	Age        string `json:"age"`
}

// formatAge formats the duration since last update as a human-readable string
func formatAge(updatedAt time.Time) string {
	duration := time.Since(updatedAt)

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return "just now"
}

func NewAccountResponse(account *entity.Account) *AccountResponse {
	return &AccountResponse{
		ID:         account.ID,
		UniqueID:   account.UniqueID,
		Name:       account.Name,
		Username:   account.Username,
		Email:      account.Email,
		Platform:   account.Platform,
		StudioID:   fmt.Sprint(account.StudioID),
		StudioName: account.Studio.Name,
		Device:     account.Device,
		Cookie:     account.Cookie,
		IsActive:   account.IsActive,
		Age:        formatAge(account.CookieUpdatedAt),
	}
}

type AccountDetailResponse struct {
	ID         uint   `json:"id"`
	UniqueID   string `json:"unique_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Platform   string `json:"platform"`
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	Cookie     string `json:"cookie"`
	Device     string `json:"device"`
	IsActive   bool   `json:"is_active"`
	Age        string `json:"age"`
}
