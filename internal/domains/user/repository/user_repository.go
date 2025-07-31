package repository

import "github.com/royhairul/live-studio-api/internal/domains/user/entity"

type UserRepository interface {
	FindAll() ([]entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	Create(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(id string) error
}
