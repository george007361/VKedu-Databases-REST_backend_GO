package service

import (
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(newUserData models.UserCreate) ([]*models.User, models.Error) {
	// Check
	users, err := s.repo.GetUserProfilesByEmailOrNickname(newUserData.Email, newUserData.Nickname)

	// Пользователи с такими данными уже есть
	if err.Code == http.StatusOK {
		return users, models.Error{Code: http.StatusConflict, Message: "User with such email or nickname already exists"}
	}

	// Какая то ошибка, кроме notfound
	if err.Code != http.StatusNotFound {
		return nil, err
	}

	// Такого юзера нет, создаём
	err = s.repo.CreateUser(newUserData)

	// Не смогли созадть, ошибка
	if err.Code != http.StatusCreated {
		return nil, err
	}

	data := make([]*models.User, 1)
	data = append(data, &models.User{
		About:    newUserData.About,
		Email:    newUserData.Email,
		FullName: newUserData.FullName,
		Nickname: newUserData.Nickname,
	})

	return data, err
}

func (s *UserService) GetUserProfile(nickname string) (models.User, models.Error) {
	return s.repo.GetUserProfile(nickname)
}

func (s *UserService) UpdateUserProfile(updatedData models.UserUpdate) (models.User, models.Error) {
	if updatedData.About == "" && updatedData.Email == "" && updatedData.FullName == "" {
		return s.repo.GetUserProfile(updatedData.Nickname)
	}
	return s.repo.UpdateUserProfile(updatedData)
}
