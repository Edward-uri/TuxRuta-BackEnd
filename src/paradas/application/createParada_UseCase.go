package application

import (
	"errors"
	"main/src/paradas/domain"
	"main/src/paradas/domain/entities"
	rutasDomain "main/src/rutas/domain"
)

type CreateParadaUseCase struct {
	ParadasRepository domain.ParadasRepository
	RutaRepository    rutasDomain.IRutasRepository
}

func NewCreateParadaUseCase(paradasRepo domain.ParadasRepository, rutaRepo rutasDomain.IRutasRepository) *CreateParadaUseCase {
	return &CreateParadaUseCase{
		ParadasRepository: paradasRepo,
		RutaRepository:    rutaRepo,
	}
}
func (uc *CreateParadaUseCase) Execute(nombre string, ubicacion entities.Ubicacion, rutaID int, creadoPor int) (*entities.Parada, error) {
	if nombre == "" {
		return nil, errors.New("nombre is required")
	}
	if rutaID <= 0 {
		return nil, errors.New("valid ruta_id is required")
	}

	ruta, err := uc.RutaRepository.GetRutaById(rutaID)
	if err != nil {
		return nil, err
	}
	if ruta == nil {
		return nil, errors.New("ruta not found")
	}

	parada := entities.NewParada(nombre, ubicacion, rutaID, creadoPor)

	paradaID, err := uc.ParadasRepository.CrearParada(*parada)
	if err != nil {
		return nil, err
	}

	parada.ID = paradaID

	return parada, nil
}
