package repository

import "github.com/royhairul/live-studio-api/internal/domains/live/entity"

type LiveRepository interface {
	FindAll() ([]*entity.Live, error)
	FindById(id string)
}
