package application

import (
	"main/src/colectivo/domain"
	"main/src/colectivo/domain/entities"
)

type GetColectivosUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewGetColectivosUseCase(colectivoRepo domain.IColectivoRepository) *GetColectivosUseCase {
	return &GetColectivosUseCase{colectivoRepository: colectivoRepo}
}

func (uc *GetColectivosUseCase) Execute() ([]entities.Colectivo, error) {
	colectivos, err := uc.colectivoRepository.GetColectivos()
	if err != nil {
		return nil, err
	}
	return colectivos, nil
}
