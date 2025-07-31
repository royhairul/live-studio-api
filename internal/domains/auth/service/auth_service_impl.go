package service

import (
	"errors"

	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/helpers/errorhandler"
	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
	"github.com/royhairul/live-studio-api/internal/domains/auth/repository"
	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	userrepo "github.com/royhairul/live-studio-api/internal/domains/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	UserRepository userrepo.UserRepository
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo userrepo.UserRepository,
) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepo,
		UserRepository: userRepo,
	}
}

func (s *AuthServiceImpl) Login(user params.LoginRequest) (params.LoginResponse, error) {
	existingUser, err := s.UserRepository.FindByEmail(user.Email)
	if err != nil {
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email atau password salah")
	}

	if existingUser == nil {
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email atau password salah")
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email atau password salah")
	}

	token, err := helpers.GenerateTokenJWT(existingUser)
	if err != nil {
		return params.LoginResponse{}, err
	}

	return params.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *AuthServiceImpl) Register(user params.RegisterRequest) (params.RegisterResponse, error) {
	// cek email
	existingUser, _ := s.UserRepository.FindByEmail(user.Email)
	if existingUser != nil {
		return params.RegisterResponse{}, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return params.RegisterResponse{}, err
	}

	newUser := &entity.User{
		Name:     user.Name,
		Email:    user.Email,
		RoleID:   1, // default role
		Password: string(hashedPassword),
	}

	createdUser, err := s.UserRepository.Create(newUser)
	if err != nil {
		return params.RegisterResponse{}, err
	}

	return params.RegisterResponse{
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  createdUser.Role.Name,
	}, nil
}

func (s *AuthServiceImpl) ForgotPassword(email params.ForgotPasswordRequest) (params.ForgotPasswordResponse, error) {
	_, err := s.UserRepository.FindByEmail(email.Email)
	if err != nil {
		return params.ForgotPasswordResponse{}, errors.New("email not registered")
	}

	otp, err := s.AuthRepository.ForgotPassword(email)

	if err != nil {
		return params.ForgotPasswordResponse{}, err
	}
	return params.ForgotPasswordResponse{
		Otp: otp.Otp,
	}, nil
}

func (s *AuthServiceImpl) ResetPassword(password params.ResetPasswordRequest) (params.ChangePasswordResponse, error) {

	_, err := s.AuthRepository.ResetPassword(password)
	if err != nil {
		return params.ChangePasswordResponse{}, err
	}
	return params.ChangePasswordResponse{
		Message: "password berhasil diubah",
	}, nil
}

func (s *AuthServiceImpl) VerifyOtp(otp params.VerifyOTPRequest) (params.ChangePasswordResponse, error) {

	_, err := s.AuthRepository.VerifyOtp(otp)
	if err != nil {
		return params.ChangePasswordResponse{}, err
	}
	return params.ChangePasswordResponse{
		Message: "otp valid",
	}, nil
}
