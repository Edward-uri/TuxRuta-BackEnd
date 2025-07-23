package application

import (
	"main/src/rutas/domain"
)

type DeleteRutaUseCase struct {
	rutaRepository domain.IRutasRepository
}

func NewDeleteRutaUseCase(rutaRepository domain.IRutasRepository) *DeleteRutaUseCase {
	return &DeleteRutaUseCase{
		rutaRepository: rutaRepository,
	}
}

func (uc *DeleteRutaUseCase) Execute(id int) error {
	return uc.rutaRepository.DeleteRuta(id)
}
