package service

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type ManagmentService struct {
	repo repository.Managment
}

func NewManagmentService(repo repository.Managment) *ManagmentService {
	return &ManagmentService{repo: repo}
}

func (s *ManagmentService) Clear() models.Error {
	return s.repo.Clear()
}

func (s *ManagmentService) GetStatus() (models.Status, models.Error) {
	return s.repo.GetStatus()
}
