package entities

import "time"

type Point struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Order int     `json:"order"`
}

type PathData struct {
	Points []Point `json:"points"`
}

type Ruta struct {
	ID            int       `json:"id"`
	Nombre        string    `json:"nombre"`
	Descripcion   *string   `json:"descripcion,omitempty"`
	PathData      PathData  `json:"path_data"`
	Activa        bool      `json:"activa"`
	CreadoPor     int       `json:"creado_por"`
	ModificadoPor *int      `json:"modificado_por,omitempty"`
	CreadoEn      time.Time `json:"creado_en"`
	ModificadoEn  time.Time `json:"modificado_en"`
}

func NewRuta(nombre, descripcion string, pathData PathData, creadoPor int) *Ruta {
	return &Ruta{
		Nombre:       nombre,
		Descripcion:  &descripcion,
		PathData:     pathData,
		Activa:       true,
		CreadoPor:    creadoPor,
		CreadoEn:     time.Now(),
		ModificadoEn: time.Now(),
	}
}
