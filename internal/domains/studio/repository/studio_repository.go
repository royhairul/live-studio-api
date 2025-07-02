package repository

import "github.com/royhairul/live-studio-api/internal/domains/studio/entity"

type StudioRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Studio, error)
	FindByID(id string) (*entity.Studio, error)
	Create(studio *entity.Studio) error
	Save(studio *entity.Studio) error
	Delete(id string) error
}
