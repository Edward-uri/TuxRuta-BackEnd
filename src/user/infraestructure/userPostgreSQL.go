package infraestructure

import (
	"database/sql"
	"main/src/user/domain/entities"
)

type UserPostgreSQL struct {
	db *sql.DB
}

func NewUserPostgreSQL(db *sql.DB) *UserPostgreSQL {
	return &UserPostgreSQL{db: db}
}
func (r *UserPostgreSQL) Create(user entities.Usuario) error {
	query := `INSERT INTO usuario (email, password, rol, activo, ultimo_acceso, creado_en) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING id`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password, // Ya viene hasheado del caso de uso
		user.Rol,
		user.Activo,
		user.UltimoAcceso,
		user.CreadoEn,
	).Scan(&user.ID)

	return err
}
func (r *UserPostgreSQL) Delete(id int) error {
	query := `DELETE FROM usuario WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserPostgreSQL) Modify(id int, user entities.Usuario) error {
	query := `UPDATE usuario SET email = $1, password = $2, rol = $3, activo = $4, ultimo_acceso = $5 
	          WHERE id = $6`
	_, err := r.db.Exec(query, user.Email, user.Password, user.Rol, user.Activo, user.UltimoAcceso, id)
	return err
}

func (r *UserPostgreSQL) GetUsers() ([]entities.Usuario, error) {
	query := `SELECT id, email, password, rol, activo, ultimo_acceso, creado_en FROM usuario`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.Usuario
	for rows.Next() {
		var user entities.Usuario
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Rol, &user.Activo, &user.UltimoAcceso, &user.CreadoEn); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserPostgreSQL) GetUserById(id int) (entities.Usuario, error) {
	query := `SELECT id, email, password, rol, activo, ultimo_acceso, creado_en FROM usuario WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user entities.Usuario
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Rol, &user.Activo, &user.UltimoAcceso, &user.CreadoEn); err != nil {
		if err == sql.ErrNoRows {
			return entities.Usuario{}, nil // No user found
		}
		return entities.Usuario{}, err
	}
	return user, nil
}
