package entities

import "time"

type Parada struct {
	ID        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	Ubicacion Ubicacion `json:"ubicacion"`
	RutaID    int       `json:"ruta_id"`
	Activa    bool      `json:"activa"`
	CreadoPor int       `json:"creado_por"`
	CreadoEn  string    `json:"creado_en"`
}

type Ubicacion struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
	Order    int     `json:"order"`
}

func NewParada(nombre string, ubicacion Ubicacion, rutaID int, creadoPor int) *Parada {
	return &Parada{
		Nombre:    nombre,
		Ubicacion: ubicacion,
		RutaID:    rutaID,
		Activa:    true,
		CreadoPor: creadoPor,
		CreadoEn:  time.Now().Format(time.RFC3339),
	}
}
