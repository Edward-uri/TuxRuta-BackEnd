package application

import (
	"errors"
	"main/src/colectivo/domain"
)

type DeleteColectivoUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewDeleteColectivoUseCase(colectivoRepo domain.IColectivoRepository) *DeleteColectivoUseCase {
	return &DeleteColectivoUseCase{colectivoRepository: colectivoRepo}
}

func (uc *DeleteColectivoUseCase) Execute(id int) error {
	if id <= 0 {
		return errors.New("invalid colectivo ID")
	}

	return uc.colectivoRepository.Delete(id)
}
