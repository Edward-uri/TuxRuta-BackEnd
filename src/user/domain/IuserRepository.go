package domain

import "main/src/user/domain/entities"

type IUserRepository interface {
	Create(user entities.Usuario) error
	Modify(id int, user entities.Usuario) error
	Delete(id int) error
	GetUsers() ([]entities.Usuario, error)
	GetUserById(id int) (entities.Usuario, error)
}
