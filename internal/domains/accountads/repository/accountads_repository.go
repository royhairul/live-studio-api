package repository

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountads/entity"
)

type AccountadsRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Accountads, error)
	FindByID(id string) (*entity.Accountads, error)
	Create(data *entity.Accountads) (*entity.Accountads, error)
	Update(data *entity.Accountads) (*entity.Accountads, error)
	Delete(id string) error

	FindByDateAndAccount(date *time.Time, AccountID string) (*entity.Accountads, error)
}
