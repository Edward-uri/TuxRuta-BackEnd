package infraestructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"main/src/rutas/domain/entities"
	"time"
)

type RutasPostgreSQL struct {
	db *sql.DB
}

func NewRutasPostgreSQL(db *sql.DB) *RutasPostgreSQL {
	return &RutasPostgreSQL{db: db}
}

func (r *RutasPostgreSQL) CreateRuta(ruta entities.Ruta) error {
	// Serializa PathData a JSON
	pathDataJSON, err := json.Marshal(ruta.PathData)
	if err != nil {
		return err
	}

	query := `INSERT INTO ruta (nombre, descripcion, path_data, activa, creado_por, creado_en, modificado_en) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.Exec(query,
		ruta.Nombre,
		ruta.Descripcion,
		pathDataJSON,
		ruta.Activa,
		ruta.CreadoPor,
		ruta.CreadoEn,
		ruta.ModificadoEn)

	return err
}

func (r *RutasPostgreSQL) DeleteRuta(id int) error {
	query := `DELETE FROM ruta WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
func (r *RutasPostgreSQL) ModifyRuta(id int, ruta entities.Ruta) error {
	query := `UPDATE ruta SET nombre = $1, descripcion = $2, path_data = $3, activa = $4, modificado_por = $5, modificado_en = $6 WHERE id = $7`
	_, err := r.db.Exec(query,
		ruta.Nombre,
		ruta.Descripcion,
		ruta.PathData,
		ruta.Activa,
		ruta.ModificadoPor,
		time.Now(),
		id,
	)
	return err
}
func (r *RutasPostgreSQL) GetRutas(ctx context.Context) ([]entities.Ruta, error) {
	query := `SELECT id, nombre, descripcion, path_data, activa, creado_por, modificado_por, creado_en, modificado_en 
              FROM ruta WHERE activa = true ORDER BY creado_en ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rutas []entities.Ruta

	for rows.Next() {
		var ruta entities.Ruta
		var pathDataJSON []byte

		err := rows.Scan(&ruta.ID, &ruta.Nombre, &ruta.Descripcion, &pathDataJSON,
			&ruta.Activa, &ruta.CreadoPor, &ruta.ModificadoPor, &ruta.CreadoEn, &ruta.ModificadoEn)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(pathDataJSON, &ruta.PathData)
		if err != nil {
			return nil, err
		}

		rutas = append(rutas, ruta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rutas, nil
}

func (r *RutasPostgreSQL) GetRutaById(id int) (*entities.Ruta, error) {
	query := `SELECT id, nombre, descripcion, path_data, activa, creado_por, modificado_por, creado_en, modificado_en 
              FROM ruta WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var ruta entities.Ruta
	var pathDataJSON []byte

	err := row.Scan(&ruta.ID, &ruta.Nombre, &ruta.Descripcion, &pathDataJSON,
		&ruta.Activa, &ruta.CreadoPor, &ruta.ModificadoPor, &ruta.CreadoEn, &ruta.ModificadoEn)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(pathDataJSON, &ruta.PathData)
	if err != nil {
		return nil, err
	}

	return &ruta, nil
}

func (r *RutasPostgreSQL) GetRutaByNombre(nombre string) (*entities.Ruta, error) {
	query := `SELECT id, nombre, descripcion, path_data, activa, creado_por, modificado_por, creado_en, modificado_en FROM ruta WHERE nombre = $1`
	row := r.db.QueryRow(query, nombre)

	var ruta entities.Ruta
	var pathDataJSON []byte
	err := row.Scan(&ruta.ID, &ruta.Nombre, &ruta.Descripcion, &pathDataJSON, &ruta.Activa,
		&ruta.CreadoPor, &ruta.ModificadoPor, &ruta.CreadoEn, &ruta.ModificadoEn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(pathDataJSON, &ruta.PathData)
	if err != nil {
		return nil, err
	}

	return &ruta, nil
}

/*
func (r *RutasPostgreSQL) GetRutaAndParadas(nombre string) (*entities.Ruta, error) {
    // Primero obtienes la ruta
    query := `SELECT id, nombre, descripcion, path_data, activa, creado_por, modificado_por, creado_en, modificado_en
              FROM ruta WHERE nombre = $1`
    row := r.db.QueryRow(query, nombre)

    var ruta entities.Ruta
    var pathDataJSON []byte
    err := row.Scan(&ruta.ID, &ruta.Nombre, &ruta.Descripcion, &pathDataJSON, &ruta.Activa,
        &ruta.CreadoPor, &ruta.ModificadoPor, &ruta.CreadoEn, &ruta.ModificadoEn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Ruta no encontrada
        }
        return nil, err
    }

    // Deserializa el path_data JSON
    err = json.Unmarshal(pathDataJSON, &ruta.PathData)
    if err != nil {
        return nil, err
    }

    // Luego obtienes las paradas de esa ruta
    paradasQuery := `SELECT id, nombre, ubicacion, datos_extra, activa, creado_por, creado_en
                     FROM parada WHERE ruta_id = $1 AND activa = true ORDER BY id`
    rows, err := r.db.Query(paradasQuery, ruta.ID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var paradas []entities.Parada // Necesitarás crear esta entidad
    for rows.Next() {
        var parada entities.Parada
        var ubicacionJSON, datosExtraJSON []byte
        err := rows.Scan(&parada.ID, &parada.Nombre, &ubicacionJSON, &datosExtraJSON,
            &parada.Activa, &parada.CreadoPor, &parada.CreadoEn)
        if err != nil {
            return nil, err
        }

        // Deserializa los campos JSON
        err = json.Unmarshal(ubicacionJSON, &parada.Ubicacion)
        if err != nil {
            return nil, err
        }

        if datosExtraJSON != nil {
            err = json.Unmarshal(datosExtraJSON, &parada.DatosExtra)
            if err != nil {
                return nil, err
            }
        }

        paradas = append(paradas, parada)
    }

    // Agrega las paradas a la ruta
    ruta.Paradas = paradas
    return &ruta, nil
}


func (r *RutasPostgreSQL) GetRutaParadaColectivo(nombre string) (*entities.Ruta, error) {} */
