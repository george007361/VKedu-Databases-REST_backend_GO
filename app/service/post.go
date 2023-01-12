package service

import (
	"net/http"

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
	postData, err := s.repo.GetPostData(id)
	if err.Code != http.StatusOK {
		return postData, err
	}
	// Empty change

	if newData.Message == "" {
		return postData, err
	}

	// Check same
	if newData.Message == postData.Message {
		return postData, err
	}

	return s.repo.UpdatePostData(newData, id)
}
