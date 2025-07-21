package application

import (
	"errors"
	"main/src/colectivo/domain"
	"main/src/colectivo/domain/entities"
)

type ModifyColectivoUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewModifyColectivoUseCase(colectivoRepo domain.IColectivoRepository) *ModifyColectivoUseCase {
	return &ModifyColectivoUseCase{colectivoRepository: colectivoRepo}
}

func (uc *ModifyColectivoUseCase) Execute(id int, colectivo *entities.Colectivo) error {
	if colectivo.Matricula == "" {
		return errors.New("please provide a valid matricula")
	}

	if colectivo.Ruta_id <= 0 {
		return errors.New("ruta_id must be greater than 0")
	}

	return uc.colectivoRepository.Modify(id, *colectivo)
}
