package service

import (
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
	"github.com/sirupsen/logrus"
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
	if len(newPostsData) < 1 {
		logrus.Println("NO POSTS")
		return []models.Post{}, models.Error{Code: http.StatusCreated}
	}
	return s.repo.CreatePostsByThreadSlug(newPostsData, slug)
}

func (s *ThreadService) CreatePostsByThreadId(newPostsData []models.Post, id int) ([]models.Post, models.Error) {
	if len(newPostsData) < 1 {
		logrus.Println("NO POSTS")
		return []models.Post{}, models.Error{Code: http.StatusCreated}
	}
	return s.repo.CreatePostsByThreadId(newPostsData, id)
}

func (s *ThreadService) UpdateThreadBySlug(newData models.UpdateThread, slug string) (models.Thread, models.Error) {
	return s.repo.UpdateThreadBySlug(newData, slug)
}
func (s *ThreadService) UpdateThreadById(newData models.UpdateThread, id int) (models.Thread, models.Error) {
	return s.repo.UpdateThreadById(newData, id)
}

func (s *ThreadService) GetThreadPostsBySlug(params models.ThreadGetPostsParams, slug string) ([]models.Post, models.Error) {
	return s.repo.GetThreadPostsBySlug(params, slug)
}
func (s *ThreadService) GetThreadPostsById(params models.ThreadGetPostsParams, id int) ([]models.Post, models.Error) {
	return s.repo.GetThreadPostsById(params, id)
}
func (s *ThreadService) VoteThreadBySlug(vote models.Vote, slug string) (models.Thread, models.Error) {
	return s.repo.VoteThreadBySlug(vote, slug)
}
func (s *ThreadService) VoteThreadById(vote models.Vote, id int) (models.Thread, models.Error) {
	return s.repo.VoteThreadById(vote, id)
}
