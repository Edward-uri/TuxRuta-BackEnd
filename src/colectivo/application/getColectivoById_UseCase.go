package application

import (
	"errors"
	"main/src/colectivo/domain"
	"main/src/colectivo/domain/entities"
)

type GetColectivoByIdUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewGetColectivoByIdUseCase(colectivoRepo domain.IColectivoRepository) *GetColectivoByIdUseCase {
	return &GetColectivoByIdUseCase{colectivoRepository: colectivoRepo}
}

func (uc *GetColectivoByIdUseCase) Execute(id int) (entities.Colectivo, error) {
	if id <= 0 {
		return entities.Colectivo{}, errors.New("invalid colectivo ID")
	}
	return uc.colectivoRepository.GetColectivoById(id)
}
