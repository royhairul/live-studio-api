package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
	"github.com/royhairul/live-studio-api/internal/domains/auth/repository"
	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	userrepo "github.com/royhairul/live-studio-api/internal/domains/user/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	helpers "github.com/royhairul/live-studio-api/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	roleservice "github.com/royhairul/live-studio-api/internal/domains/role/service"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	UserRepository userrepo.UserRepository
	RoleService    roleservice.RoleService
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo userrepo.UserRepository,
	roleSvc roleservice.RoleService,
) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepo,
		UserRepository: userRepo,
		RoleService:    roleSvc,
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

func (s *AuthServiceImpl) ForgotPassword(input params.ForgotPasswordRequest) (email string, err error) {
	_, err = s.UserRepository.FindByEmail(input.Email)
	if err != nil {
		return " ", errors.New("email not registered")
	}
	otp, err := s.AuthRepository.ForgotPassword(input)

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")

	port, _ := strconv.Atoi(portStr)

	d := gomail.NewDialer(host, port, from, password)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", input.Email)
	m.SetHeader("Subject", "Reset Password OTP")
	body := fmt.Sprintf("Hello, <br>Your OTP is: <b>%s</b>", otp.Otp)
	m.SetBody("text/html", body)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	if err != nil {
		return " ", err
	}
	return input.Email, nil
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

// Me implements AuthService.
func (s *AuthServiceImpl) Me(userId string) (params.MeResponse, error) {
	user, err := s.UserRepository.FindByID(userId)
	if err != nil {
		return params.MeResponse{}, fmt.Errorf("failed to find user with ID %s: %w", userId, err)
	}

	log.Println("Role found:", user.RoleID)
	log.Println("Role found:", fmt.Sprintf("%d", user.RoleID))

	role, err := s.RoleService.FindByID(fmt.Sprintf("%d", user.RoleID))
	if err != nil {
		return params.MeResponse{}, fmt.Errorf("failed to find role with ID %d: %w", user.RoleID, err)
	}

	permissions := []string{}
	for _, p := range role.Permissions {
		permissions = append(permissions, p.Name)
	}

	return params.MeResponse{
		Name:        user.Name,
		Role:        role.Name,
		Permissions: permissions,
	}, nil
}
