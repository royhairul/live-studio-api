package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
	"github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
)

type AccountsessionRepository interface {
	// TODO: define repository methods
	FindAll(filter params.AccountsessionFilter) ([]*entity.Accountsession, error)
	FindOne(filter params.AccountsessionFilter) (*entity.Accountsession, error)
	Create(data *entity.Accountsession) (*entity.Accountsession, error)
	Update(data *entity.Accountsession) (*entity.Accountsession, error)
	Delete(id string) error

	FindAllByAttendanceID(id string) ([]*entity.Accountsession, error)
	FindAllByAccountID(id string) ([]*entity.Accountsession, error)
	FindAllByStudioID(id string) ([]*entity.Accountsession, error)
}
