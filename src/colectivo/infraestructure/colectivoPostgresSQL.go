package infraestructure

import (
	"database/sql"
	"main/src/colectivo/domain/entities"
	"time"
)

type ColectivoPostgreSQL struct {
	db *sql.DB
}

func NewColectivoPostgreSQL(db *sql.DB) *ColectivoPostgreSQL {
	return &ColectivoPostgreSQL{db: db}
}

func (r *ColectivoPostgreSQL) Create(colectivo entities.Colectivo) error {
	query := `INSERT INTO colectivo (matricula, ruta_id, activo, creado_por, creado_en)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id`
	err := r.db.QueryRow(
		query,
		colectivo.Matricula,
		colectivo.Ruta_id,
		colectivo.Activo,
		colectivo.CreadoPor, // <-- ID del usuario creador
		time.Now(),
	).Scan(&colectivo.ID)
	return err
}

func (r *ColectivoPostgreSQL) Delete(id int) error {
	query := `DELETE FROM colectivo WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ColectivoPostgreSQL) Modify(id int, colectivo entities.Colectivo) error {
	query := `UPDATE colectivo SET matricula = $1, ruta_id = $2, activo = $3 WHERE id = $4`
	_, err := r.db.Exec(query,
		colectivo.Matricula,
		colectivo.Ruta_id,
		colectivo.Activo,
		id,
	)
	return err
}

func (r *ColectivoPostgreSQL) GetColectivos() ([]entities.Colectivo, error) {
	query := `SELECT id, matricula, ruta_id, activo, creado_por, creado_en FROM colectivo`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colectivos []entities.Colectivo
	for rows.Next() {
		var c entities.Colectivo
		var creadoEn time.Time
		if err := rows.Scan(&c.ID, &c.Matricula, &c.Ruta_id, &c.Activo, &c.CreadoPor, &creadoEn); err != nil {
			return nil, err
		}
		c.CreateAt = creadoEn
		colectivos = append(colectivos, c)
	}
	return colectivos, nil
}

func (r *ColectivoPostgreSQL) GetColectivoById(id int) (entities.Colectivo, error) {
	query := `SELECT id, matricula, ruta_id, activo, creado_por, creado_en FROM colectivo WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var c entities.Colectivo
	var creadoEn time.Time
	if err := row.Scan(&c.ID, &c.Matricula, &c.Ruta_id, &c.Activo, &c.CreadoPor, &creadoEn); err != nil {
		if err == sql.ErrNoRows {
			return entities.Colectivo{}, nil
		}
		return entities.Colectivo{}, err
	}
	c.CreateAt = creadoEn
	return c, nil
}

func (r *ColectivoPostgreSQL) GetColectivoByMatricula(matricula string) (*entities.Colectivo, error) {
	query := `SELECT id, matricula, ruta_id, activo, creado_por, creado_en FROM colectivo WHERE matricula = $1`
	row := r.db.QueryRow(query, matricula)

	var c entities.Colectivo
	var creadoEn time.Time
	if err := row.Scan(&c.ID, &c.Matricula, &c.Ruta_id, &c.Activo, &c.CreadoPor, &creadoEn); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.CreateAt = creadoEn
	return &c, nil
}
