package service

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) models.Error {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserProfile(nickname string) (models.User, models.Error) {
	return s.repo.GetUserProfile(nickname)
}

func (s *UserService) UpdateUserProfile(userData models.User) (models.User, models.Error) {
	return s.repo.UpdateUserProfile(userData)
}
