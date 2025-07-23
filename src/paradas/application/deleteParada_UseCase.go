package application

import (
	"errors"
	"main/src/paradas/domain"
)

type DeleteParadaUseCase struct {
	ParadasRepository domain.ParadasRepository
}

func NewDeleteParadaUseCase(repo domain.ParadasRepository) *DeleteParadaUseCase {
	return &DeleteParadaUseCase{
		ParadasRepository: repo,
	}
}

func (uc *DeleteParadaUseCase) Execute(id int) error {
	if id <= 0 {
		return errors.New("ID is required")
	}
	uc.ParadasRepository.ObtenerParadaById(id)
	if err := uc.ParadasRepository.EliminarParada(id); err != nil {
		return errors.New("Error deleting parada: " + err.Error())
	}
	return nil
}
