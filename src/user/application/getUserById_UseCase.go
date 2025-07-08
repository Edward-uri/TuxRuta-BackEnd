package application

import (
	"main/src/user/domain"
	"main/src/user/domain/entities"
)

type GetUserByIdUseCase struct {
	repository domain.IUserRepository
}

func NewGetUserByIdUseCase(repo domain.IUserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{
		repository: repo,
	}
}

func (uc *GetUserByIdUseCase) Execute(id int) (entities.Usuario, error) {
	return uc.repository.GetUserById(id)
}
