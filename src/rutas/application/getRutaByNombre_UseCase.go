package application

import (
	"main/src/rutas/domain"
	"main/src/rutas/domain/entities"
)

type GetRutaByNombreUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewGetRutaByNombreUseCase(rutaRepo domain.IRutasRepository) *GetRutaByNombreUseCase {
	return &GetRutaByNombreUseCase{rutaRepository: rutaRepo}
}

func (uc *GetRutaByNombreUseCase) Execute(nombre string) (*entities.Ruta, error) {
	ruta, err := uc.rutaRepository.GetRutaByNombre(nombre)
	if err != nil {
		return nil, err
	}
	return ruta, nil
}
