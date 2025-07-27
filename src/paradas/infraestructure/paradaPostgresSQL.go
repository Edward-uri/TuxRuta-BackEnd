package infraestructure

import (
	"database/sql"
	"encoding/json"
	"main/src/paradas/domain/entities"
	"time"
)

type ParadasPostgreSQL struct {
	db *sql.DB
}

func NewParadasPostgreSQL(db *sql.DB) *ParadasPostgreSQL {
	return &ParadasPostgreSQL{db: db}
}

// CrearParada: retorna el ID insertado
func (r *ParadasPostgreSQL) CrearParada(parada entities.Parada) (int, error) {
	ubicacionJSON, err := json.Marshal(parada.Ubicacion)
	if err != nil {
		return 0, err
	}
	query := `INSERT INTO parada (nombre, ubicacion, ruta_id, activa, creado_por, creado_en)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err = r.db.QueryRow(query, parada.Nombre, ubicacionJSON, parada.RutaID, parada.Activa, parada.CreadoPor, time.Now()).Scan(&id)
	return id, err
}

// ActualizarParada
func (r *ParadasPostgreSQL) ActualizarParada(parada entities.Parada) error {
	ubicacionJSON, err := json.Marshal(parada.Ubicacion)
	if err != nil {
		return err
	}
	query := `UPDATE parada SET nombre = $1, ubicacion = $2, ruta_id = $3 WHERE id = $4`
	_, err = r.db.Exec(query, parada.Nombre, ubicacionJSON, parada.RutaID, parada.ID)
	return err
}

// ObtenerParadaPorID
func (r *ParadasPostgreSQL) ObtenerParadaById(id int) (*entities.Parada, error) {
	query := `SELECT id, nombre, ubicacion, ruta_id, activa, creado_por, creado_en FROM parada WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var parada entities.Parada
	var ubicacionJSON []byte
	err := row.Scan(&parada.ID, &parada.Nombre, &ubicacionJSON, &parada.RutaID, &parada.Activa, &parada.CreadoPor, &parada.CreadoEn)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(ubicacionJSON, &parada.Ubicacion)
	if err != nil {
		return nil, err
	}
	return &parada, nil
}

// ObtenerParadasPorRuta
func (r *ParadasPostgreSQL) ObtenerParadasPorRuta(rutaID int) ([]entities.Parada, error) {
	query := `SELECT id, nombre, ubicacion, ruta_id, activa, creado_por, creado_en FROM parada WHERE ruta_id = $1 AND activa = true ORDER BY id ASC`
	rows, err := r.db.Query(query, rutaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paradas []entities.Parada
	for rows.Next() {
		var parada entities.Parada
		var ubicacionJSON []byte
		err := rows.Scan(&parada.ID, &parada.Nombre, &ubicacionJSON, &parada.RutaID, &parada.Activa, &parada.CreadoPor, &parada.CreadoEn)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(ubicacionJSON, &parada.Ubicacion)
		if err != nil {
			return nil, err
		}
		paradas = append(paradas, parada)
	}
	return paradas, nil
}

// EliminarParada
func (r *ParadasPostgreSQL) EliminarParada(id int) error {
	query := `DELETE FROM parada WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ParadasPostgreSQL) ObtenerParadas() ([]entities.Parada, error) {
	query := `SELECT id, nombre, ubicacion, ruta_id, activa, creado_por, creado_en FROM parada ORDER BY id ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paradas []entities.Parada
	for rows.Next() {
		var parada entities.Parada
		var ubicacionJSON []byte
		err := rows.Scan(&parada.ID, &parada.Nombre, &ubicacionJSON, &parada.RutaID, &parada.Activa, &parada.CreadoPor, &parada.CreadoEn)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(ubicacionJSON, &parada.Ubicacion)
		if err != nil {
			return nil, err
		}
		paradas = append(paradas, parada)
	}
	return paradas, nil
}
