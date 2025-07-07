package application

import (
	"errors"
	"main/src/user/application/services"
	"main/src/user/domain"
	"main/src/user/domain/entities"
	"time"
)

type CreateUserUseCase struct {
	userRepository  domain.IUserRepository
	passwordService services.PasswordService
}

func NewCreateUserUseCase(userRepo domain.IUserRepository, passwordService services.PasswordService) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepo, passwordService: passwordService}
}

func (uc *CreateUserUseCase) Execute(user *entities.Usuario) error {
	// Validaciones de negocio
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Hashear password usando el servicio
	hashedPassword, err := uc.passwordService.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Preparar usuario para persistencia
	user.Password = hashedPassword
	user.Activo = true
	user.CreadoEn = time.Now()

	// Guardar en repositorio
	return uc.userRepository.Create(*user)
}
