package application

import (
	"main/src/paradas/domain"
	"main/src/paradas/domain/entities"
)

type GetParadaUseCase struct {
	paradaRepository domain.ParadasRepository
}

func NewGetParadaUseCase(repo domain.ParadasRepository) *GetParadaUseCase {
	return &GetParadaUseCase{
		paradaRepository: repo,
	}
}

func (uc *GetParadaUseCase) Execute(id int) (entities.Parada, error) {
	parada, err := uc.paradaRepository.ObtenerParadaById(id)
	if err != nil {
		return entities.Parada{}, err
	}
	return *parada, nil
}
