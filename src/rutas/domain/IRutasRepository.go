package domain

import (
	"context"
	"main/src/rutas/domain/entities"
)

type IRutasRepository interface {
	CreateRuta(ruta entities.Ruta) error
	ModifyRuta(id int, ruta entities.Ruta) error
	DeleteRuta(id int) error
	GetRutas(ctx context.Context) ([]entities.Ruta, error)
	GetRutaById(id int) (*entities.Ruta, error)
	GetRutaByNombre(nombre string) (*entities.Ruta, error)
}
