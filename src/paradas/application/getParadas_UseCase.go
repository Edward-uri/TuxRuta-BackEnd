package application

import (
	"main/src/paradas/domain"
	"main/src/paradas/domain/entities"
)

type GetParadasAllUseCase struct {
	paradaRepository domain.ParadasRepository
}

func NewGetParadasAllUseCase(paradaRepository domain.ParadasRepository) *GetParadasAllUseCase {
	return &GetParadasAllUseCase{
		paradaRepository: paradaRepository,
	}
}

func (uc *GetParadasAllUseCase) Execute() ([]entities.Parada, error) {
	paradas, err := uc.paradaRepository.ObtenerParadas()
	if err != nil {
		return nil, err
	}
	return paradas, nil
}
