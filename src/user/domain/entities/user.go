package entities

import (
	"time"
)

type Usuario struct {
	ID           int        `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"password" db:"password"`
	Rol          string     `json:"rol" db:"rol"`
	Activo       bool       `json:"activo" db:"activo"`
	UltimoAcceso *time.Time `json:"ultimo_acceso,omitempty" db:"ultimo_acceso"`
	CreadoEn     time.Time  `json:"creado_en" db:"creado_en"`
}

func NewUsuario(email, password, rol string) *Usuario {
	return &Usuario{
		Email:        email,
		Password:     password,
		Rol:          rol,
		Activo:       true,
		UltimoAcceso: nil,
		CreadoEn:     time.Now(),
	}
}
