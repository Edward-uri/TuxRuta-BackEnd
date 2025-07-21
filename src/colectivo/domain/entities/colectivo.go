package entities

import (
	"time"
)

type Colectivo struct {
	ID        int       `json:"id"`
	Matricula string    `json:"matricula"`
	Ruta_id   int       `json:"ruta_id"`
	Activo    bool      `json:"activo"`
	CreadoPor int       `json:"creado_por"`
	CreateAt  time.Time `json:"creado_en"`
}

func NewColectivo(matricula string, ruta_id int, creadoPor int) *Colectivo {
	return &Colectivo{
		Matricula: matricula,
		Ruta_id:   ruta_id,
		Activo:    true,
		CreadoPor: creadoPor,
		CreateAt:  time.Now(),
	}
}
