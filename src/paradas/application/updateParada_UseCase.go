package application

import (
	"errors"
	"main/src/paradas/domain"
	"main/src/paradas/domain/entities"
)

type UpdateParadaUseCase struct {
	ParadasRepository domain.ParadasRepository
}

func NewUpdateParadaUseCase(repo domain.ParadasRepository) *UpdateParadaUseCase {
	return &UpdateParadaUseCase{
		ParadasRepository: repo,
	}
}

func (uc *UpdateParadaUseCase) Execute(parada entities.Parada) error {
	if parada.ID <= 0 {
		return errors.New("ID is required")
	}
	if parada.RutaID < 0 {
		return errors.New("RutaID is required")
	}
	uc.ParadasRepository.ObtenerParadaById(parada.ID)
	if err := uc.ParadasRepository.ActualizarParada(parada); err != nil {
		return errors.New("Error updating parada: " + err.Error())
	}
	return uc.ParadasRepository.ActualizarParada(parada)
}
