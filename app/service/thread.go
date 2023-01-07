package service

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type ThreadService struct {
	repo repository.Thread
}

func NewThreadService(repo repository.Thread) *ThreadService {
	return &ThreadService{repo: repo}
}

func (s *ThreadService) GetThreadData(slug string) (models.Thread, models.Error) {
	return s.repo.GetThreadData(slug)
}

func (s *ThreadService) GetThreadDataById(id int) (models.Thread, models.Error) {
	return s.repo.GetThreadDataById(id)
}

func (s *ThreadService) CreatePostsByThreadSlug(newPostsData []models.Post, slug string) ([]models.Post, models.Error) {
	return s.repo.CreatePostsByThreadSlug(newPostsData, slug)
}
