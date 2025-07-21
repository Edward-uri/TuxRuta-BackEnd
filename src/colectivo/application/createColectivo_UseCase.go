package application

import (
	"errors"
	"main/src/colectivo/domain"
	"main/src/colectivo/domain/entities"
)

type CreateColectivoUseCase struct {
	colectivoRepository domain.IColectivoRepository
}

func NewCreateColectivoUseCase(colectivoRepo domain.IColectivoRepository) *CreateColectivoUseCase {
	return &CreateColectivoUseCase{colectivoRepository: colectivoRepo}
}

// Ahora recibe el id del usuario logueado
func (uc *CreateColectivoUseCase) Execute(colectivo *entities.Colectivo, creadoPor int) error {
	if colectivo.Matricula == "" {
		return errors.New("matricula is required")
	}

	if colectivo.Ruta_id <= 0 {
		return errors.New("ruta_id must be a positive integer")
	}

	colectivo.CreadoPor = creadoPor
	return uc.colectivoRepository.Create(*colectivo)
}
