package application

import (
	"errors"
	"main/src/rutas/domain"
	"main/src/rutas/domain/entities"
)

type CreateRutaUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewCreateRutaUseCase(rutaRepo domain.IRutasRepository) *CreateRutaUseCase {
	return &CreateRutaUseCase{rutaRepository: rutaRepo}
}

func (uc *CreateRutaUseCase) Execute(nombre, descripcion string, pathData entities.PathData, creadoPor int) (*entities.Ruta, error) {
	// Validaciones de negocio
	if nombre == "" {
		return nil, errors.New("nombre is required")
	}

	if len(pathData.Points) == 0 {
		return nil, errors.New("at least one point is required")
	}

	// Crear la entidad usando el constructor
	ruta := entities.NewRuta(nombre, descripcion, pathData, creadoPor)

	// Persistir en el repositorio
	err := uc.rutaRepository.CreateRuta(*ruta)
	if err != nil {
		return nil, err
	}

	return ruta, nil
}
