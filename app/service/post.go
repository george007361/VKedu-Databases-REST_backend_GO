package service

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPostData(id int) (models.Post, models.Error) {
	return s.repo.GetPostData(id)
}

func (s *PostService) UpdatePostData(newData models.PostUpdate, id int) (models.Post, models.Error) {
	return s.repo.UpdatePostData(newData, id)
}
