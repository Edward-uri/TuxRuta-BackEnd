package application

import (
	"main/src/user/domain"
	"main/src/user/domain/entities"
)

type GetAllUserUseCase struct {
	repository domain.IUserRepository
}

func NewGetAllUserUseCase(repository domain.IUserRepository) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		repository: repository,
	}
}

func (uc *GetAllUserUseCase) Execute() ([]entities.Usuario, error) {
	return uc.repository.GetUsers()
}
