package application

import "main/src/user/domain"

type DeleteUserUseCase struct {
	userRepository domain.IUserRepository
}

func NewDeleteUserUseCase(userRepo domain.IUserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepository: userRepo}
}

func (uc *DeleteUserUseCase) Execute(userId int) error {
	return uc.userRepository.Delete(userId)
}
