package repository

import (
	"database/sql"
	"errors"
	"privateCabin/entity"

	"github.com/lib/pq"
)

type UserRepository interface {
	GetUserByLogin(login, password string) (*entity.UserPublicDTO, error)
	CreateUserByData(login, password string) (*entity.UserPublicDTO, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByLogin(login, password string) (*entity.UserPublicDTO, error) {
	row := r.db.QueryRow("SELECT id, login FROM users WHERE login = $1 and password = $2", login, password)
	user := &entity.UserPublicDTO{}
	err := row.Scan(&user.ID, &user.Login)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUserByData(login, password string) (*entity.UserPublicDTO, error) {
	var user entity.UserPublicDTO
	query := `INSERT INTO users (login,password) VALUES ($1, $2) RETURNING id, login;`

	err := r.db.QueryRow(query, login, password).Scan(&user.ID, &user.Login)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return nil, errors.New("login_already_exists")
			}

		}
		return nil, err
	}
	return &user, nil

}
