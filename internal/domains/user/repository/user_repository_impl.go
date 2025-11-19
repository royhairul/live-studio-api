package repository

import (
	"errors"

	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Preload(clause.Associations).Find(&users, "role_id != ?", 1).Error
	return users, err
}

func (r *UserRepositoryImpl) FindByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload(clause.Associations).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload(clause.Associations).First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Create(user *entity.User) (*entity.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	var createdUser entity.User
	if err := r.db.Preload(clause.Associations).First(&createdUser, user.ID).Error; err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepositoryImpl) Update(user *entity.User) (*entity.User, error) {
	err := r.db.Select("Name", "Email", "RoleID", "UpdatedAt").Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}
