package application

import (
	"errors"
	"main/src/rutas/domain"
	"main/src/rutas/domain/entities"
)

type GetRutasByIdUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewGetRutasByIdUseCase(rutaRepo domain.IRutasRepository) *GetRutasByIdUseCase {
	return &GetRutasByIdUseCase{rutaRepository: rutaRepo}
}

func (uc *GetRutasByIdUseCase) Execute(id int) (*entities.Ruta, error) {
	if id <= 0 {
		return nil, errors.New("invalid ruta ID")
	}

	ruta, err := uc.rutaRepository.GetRutaById(id)
	if err != nil {
		return nil, err
	}

	if ruta == nil {
		return nil, errors.New("ruta not found")
	}

	return ruta, nil
}
