package application

import (
	"errors"
	"main/src/colectivo/domain"
	"main/src/colectivo/domain/entities"
)

type GetColectivoByMatriculaUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewGetColectivoByMatriculaUseCase(colectivoRepo domain.IColectivoRepository) *GetColectivoByMatriculaUseCase {
	return &GetColectivoByMatriculaUseCase{colectivoRepository: colectivoRepo}
}

func (uc *GetColectivoByMatriculaUseCase) Execute(matricula string) (*entities.Colectivo, error) {
	if matricula == "" {
		return nil, errors.New("matricula cannot be empty")
	}

	colectivo, err := uc.colectivoRepository.GetColectivoByMatricula(matricula)
	if err != nil {
		return nil, err
	}

	return colectivo, nil
}
