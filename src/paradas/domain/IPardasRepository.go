package domain

import (
	"main/src/paradas/domain/entities"
)

type ParadasRepository interface {
	CrearParada(parada entities.Parada) (int, error)
	ObtenerParadaById(id int) (*entities.Parada, error)
	ActualizarParada(parada entities.Parada) error
	EliminarParada(id int) error
	ObtenerParadasPorRuta(rutaID int) ([]entities.Parada, error)
	ObtenerParadas() ([]entities.Parada, error)
}
