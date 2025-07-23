package application

import (
	"errors"
	"main/src/paradas/domain"
	"main/src/paradas/domain/entities"
	rutasDomain "main/src/rutas/domain"
)

type GetParadasAndRutaUseCase struct {
	paradasRepo domain.ParadasRepository
	rutaRepo    rutasDomain.IRutasRepository
}

func NewGetParadasAndRutaUseCase(paradasRepo domain.ParadasRepository, rutaRepo rutasDomain.IRutasRepository) *GetParadasAndRutaUseCase {
	return &GetParadasAndRutaUseCase{
		paradasRepo: paradasRepo,
		rutaRepo:    rutaRepo,
	}
}

func (uc *GetParadasAndRutaUseCase) Execute(rutaID int) ([]entities.Parada, error) {
	if rutaID <= 0 {
		return nil, errors.New("invalid ruta ID")
	}

	ruta, err := uc.rutaRepo.GetRutaById(rutaID)
	if err != nil {
		return nil, err
	}
	if ruta == nil {
		return nil, errors.New("ruta not found")
	}

	paradas, err := uc.paradasRepo.ObtenerParadasPorRuta(rutaID)
	if err != nil {
		return nil, err
	}

	return paradas, nil
}
