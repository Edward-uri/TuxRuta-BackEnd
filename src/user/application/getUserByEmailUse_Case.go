package application

import (
	"errors"
	"main/src/user/domain"
	"main/src/user/domain/entities"
)

type GetUserByEmailUseCase struct {
	userRepository domain.IUserRepository
}

func NewGetUserByEmailUseCase(userRepository domain.IUserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserByEmailUseCase) Execute(email string) (entities.Usuario, error) {
	if email == "" {
		return entities.Usuario{}, errors.New("email is required")
	}

	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return entities.Usuario{}, errors.New("database error")
	}

	if user == nil {
		return entities.Usuario{}, errors.New("user not found")
	}

	return *user, nil
}
