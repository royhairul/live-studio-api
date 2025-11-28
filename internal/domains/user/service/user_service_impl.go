package service

import (
	"errors"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	"github.com/royhairul/live-studio-api/internal/domains/user/params"
	"github.com/royhairul/live-studio-api/internal/domains/user/repository"
	"golang.org/x/crypto/bcrypt"

	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
)

type userServiceImpl struct {
	repo        repository.UserRepository
	hostService hostservice.HostService
}

func NewUserService(repo repository.UserRepository, hostService hostservice.HostService) UserService {
	return &userServiceImpl{repo, hostService}
}

func (s *userServiceImpl) GetAll() ([]params.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []params.UserResponse
	for _, u := range users {
		responses = append(responses, params.UserResponse{
			ID:    u.ID.String(),
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role.Name,
		})
	}
	return responses, nil
}

func (s *userServiceImpl) GetByID(id string) (*params.UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &params.UserResponse{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role.Name,
	}, nil
}

func (s *userServiceImpl) Create(input params.CreateUserRequest) (*entity.User, error) {
	emailRegistered, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if emailRegistered != nil {
		return nil, errors.New("email already used")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		RoleID:   input.RoleID,
		Password: string(hashed),
	}

	user, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Update(id string, input params.UpdateUserRequest) (*entity.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if input.Email != u.Email {
		existing, err := s.repo.FindByEmail(input.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != u.ID {
			return nil, errors.New("email already used")
		}
		u.Email = input.Email
	}

	if input.RoleID != nil && *input.RoleID != u.RoleID {
		u.RoleID = *input.RoleID
	}

	if input.Name != u.Name {
		u.Name = input.Name
	}

	u.UpdatedAt = time.Now()

	user, err := s.repo.Update(u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Delete(id string) error {
	return s.repo.Delete(id)
}
