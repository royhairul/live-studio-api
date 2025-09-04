package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
)

type AccountsessionRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Accountsession, error)
	FindByID(id string) (*entity.Accountsession, error)
	Create(data *entity.Accountsession) (*entity.Accountsession, error)
	Update(data *entity.Accountsession) (*entity.Accountsession, error)
	Delete(id string) error

	FindAllByAttendanceID(id string) ([]*entity.Accountsession, error)
	FindAllByAccountID(id string) ([]*entity.Accountsession, error)
	FindAllByStudioID(id string) ([]*entity.Accountsession, error)
}
