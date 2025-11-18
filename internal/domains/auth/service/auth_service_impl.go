package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/royhairul/live-studio-api/internal/domains/auth/params"
	"github.com/royhairul/live-studio-api/internal/domains/auth/repository"
	"github.com/royhairul/live-studio-api/internal/domains/user/entity"
	userrepo "github.com/royhairul/live-studio-api/internal/domains/user/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
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
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email or password invalid")
	}

	if existingUser == nil {
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email or password invalid")
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return params.LoginResponse{}, errorhandler.NewNotFoundError("email or password invalid")
	}

	token, err := utils.GenerateTokenJWT(existingUser)
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
		TenantBase: tenantdb.TenantBase{
			TenantID: uuid.New().String(),
		},
	}

	createdUser, err := s.UserRepository.Create(newUser)
	if err != nil {
		return params.RegisterResponse{}, err
	}

	_, err = CreateTenantRoles(createdUser.TenantID)
	if err != nil {
		return params.RegisterResponse{}, err
	}

	return params.RegisterResponse{
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  createdUser.Role.Name,
	}, nil
}

func (s *AuthServiceImpl) ForgotPassword(input params.ForgotPasswordRequest) (string, error) {
	user, err := s.UserRepository.FindByEmail(input.Email)
	if err != nil || user == nil {
		return "", errors.New("email not registered")
	}

	otp, err := s.AuthRepository.ForgotPassword(input)
	if err != nil || otp.Otp == "" {
		return "", errors.New("failed to generate OTP")
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("invalid email port: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", input.Email)
	m.SetHeader("Subject", "Reset Password OTP")
	body := fmt.Sprintf("Hello <b>%s</b>, <br>Your OTP is: <b>%s</b>", user.Name, otp.Otp)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, from, password)
	if err := d.DialAndSend(m); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
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
func (s *AuthServiceImpl) Me(ctx context.Context) (params.MeResponse, error) {
	claims := ctx.Value("claims").(*utils.JWTPayloadDTO)

	user, err := s.UserRepository.FindByID(fmt.Sprint(claims.ID))
	if err != nil {
		return params.MeResponse{}, err
	}

	role, err := s.RoleService.FindByID(fmt.Sprint(user.RoleID))
	if err != nil {
		return params.MeResponse{}, err
	}

	permissions := make([]string, len(role.Permissions))
	for i, p := range role.Permissions {
		permissions[i] = p.Name
	}

	return params.MeResponse{
		Name:        claims.Name,
		Role:        role.Name,
		TenantID:    claims.TenantBase.TenantID,
		Permissions: permissions,
	}, nil
}
