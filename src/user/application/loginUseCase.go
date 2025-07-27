package application

import (
	"errors"
	"log"
	"main/src/user/application/services"
	"main/src/user/domain"
	"main/src/user/domain/entities"
	"time"
)

type LoginUseCase struct {
	userRepository  domain.IUserRepository
	passwordService services.PasswordService
	jwtService      services.JWTService
}

func NewLoginUseCase(
	userRepository domain.IUserRepository,
	passwordService services.PasswordService,
	jwtService services.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository:  userRepository,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  entities.Usuario `json:"user"`
	Token string           `json:"token"`
}

func (uc *LoginUseCase) Execute(request LoginRequest) (*LoginResponse, error) {
	if request.Email == "" {
		return nil, errors.New("email is required")
	}
	if request.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := uc.userRepository.GetUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.Activo {
		return nil, errors.New("user account is disabled")
	}

	err = uc.passwordService.VerifyPassword(user.Password, request.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	err = uc.userRepository.UpdateLastAccess(user.ID, now)
	if err != nil {
		log.Printf("Warning: Error updating last access for user %d: %v", user.ID, err)
	}

	user.UltimoAcceso = &now

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, user.Rol)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	user.Password = ""

	return &LoginResponse{
		User:  *user,
		Token: token,
	}, nil
}
