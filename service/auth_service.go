package service

import (
	"privateCabin/entity"
	"privateCabin/repository"
)

type UserService interface {
	GetUser(login, password string) (*entity.UserPublicDTO, error)
	CreateUser(login, password string) (*entity.UserPublicDTO, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetUser(login, password string) (*entity.UserPublicDTO, error) {
	return s.repo.GetUserByLogin(login, password)
}

func (s *userService) CreateUser(login, password string) (*entity.UserPublicDTO, error) {
	return s.repo.CreateUserByData(login, password)
}
