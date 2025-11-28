package application

import (
	"context"
	"errors"
	"main/src/rutas/domain"
	"main/src/rutas/domain/entities"
)

type GetRutasUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewGetRutasUseCase(rutaRepo domain.IRutasRepository) *GetRutasUseCase {
	return &GetRutasUseCase{rutaRepository: rutaRepo}
}

func (uc *GetRutasUseCase) Execute(ctx context.Context) ([]entities.Ruta, error) {
	rutas, err := uc.rutaRepository.GetRutas(ctx)
	if err != nil {
		return nil, err
	}
	if len(rutas) == 0 {
		return nil, errors.New("no rutas found")
	}
	return rutas, nil
}
