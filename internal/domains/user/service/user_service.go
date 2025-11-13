package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	"github.com/royhairul/live-studio-api/internal/domains/user/params"
)

type UserService interface {
	GetAll() ([]params.UserResponse, error)
	GetByID(id string) (*params.UserResponse, error)
	Create(input params.CreateUserRequest) (*entity.User, error)
	Update(id string, input params.UpdateUserRequest) (*entity.User, error)
	Delete(id string) error
}
