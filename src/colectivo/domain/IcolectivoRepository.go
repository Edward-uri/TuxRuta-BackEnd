package domain

import "main/src/colectivo/domain/entities"

type IColectivoRepository interface {
	Create(colectivo entities.Colectivo) error
	Modify(id int, colectivo entities.Colectivo) error
	Delete(id int) error
	GetColectivos() ([]entities.Colectivo, error)
	GetColectivoById(id int) (entities.Colectivo, error)
	GetColectivoByMatricula(matricula string) (*entities.Colectivo, error)
}
