package application

import (
	"errors"
	"main/src/rutas/domain"
	"main/src/rutas/domain/entities"
	"time"
)

type ModifyRutaUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewModifyRutaUseCase(rutaRepo domain.IRutasRepository) *ModifyRutaUseCase {
	return &ModifyRutaUseCase{rutaRepository: rutaRepo}
}

func (uc *ModifyRutaUseCase) Execute(id int, updatedRuta *entities.Ruta) error {
	existingRuta, err := uc.rutaRepository.GetRutaById(id)
	if err != nil {
		return err
	}

	// Verificar si la ruta existe
	if existingRuta == nil {
		return errors.New("ruta not found")
	}

	// Validaciones de negocio
	if updatedRuta.Nombre == "" {
		return errors.New("please provide a valid name")
	}
	if updatedRuta.Descripcion == nil {
		return errors.New("please provide a valid description")
	}
	if len(updatedRuta.PathData.Points) == 0 {
		return errors.New("please provide valid path data with at least one point")
	}

	// Actualizar los campos
	existingRuta.Nombre = updatedRuta.Nombre
	existingRuta.Descripcion = updatedRuta.Descripcion
	existingRuta.PathData = updatedRuta.PathData
	existingRuta.ModificadoEn = time.Now() // Actualizar timestamp

	// Pasar el struct dereferenciado al repositorio
	return uc.rutaRepository.ModifyRuta(id, *existingRuta)
}
