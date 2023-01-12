package service

import (
	"fmt"
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
	// Given : [parent, author, message]
	// Thread Slug

	// Check thread exists
	// Check forum exists
	threadData, err := s.repo.GetThreadData(slug)
	if err.Code == http.StatusNotFound {
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, slug)}
	}

	// Check empty posts list
	if len(newPostsData) < 1 {
		logrus.Println("NO POSTS")
		return []models.Post{}, models.Error{Code: http.StatusCreated}
	}

	return s.repo.CreatePosts(newPostsData, threadData.ID, threadData.Forum)
}

func (s *ThreadService) CreatePostsByThreadId(newPostsData []models.Post, id int) ([]models.Post, models.Error) {
	// Given : [parent, author, message]
	// Thread Id

	// Check thread exists
	// Check forum exists
	threadData, err := s.repo.GetThreadDataById(id)
	if err.Code == http.StatusNotFound {
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, id)}
	}

	// Check empty posts list
	if len(newPostsData) < 1 {
		logrus.Println("NO POSTS")
		return []models.Post{}, models.Error{Code: http.StatusCreated}
	}

	return s.repo.CreatePosts(newPostsData, threadData.ID, threadData.Forum)
}

func (s *ThreadService) UpdateThreadBySlug(newData models.UpdateThread, slug string) (models.Thread, models.Error) {
	// Check empty
	if newData.Message == "" && newData.Title == "" {
		return s.repo.GetThreadData(slug)
	}

	return s.repo.UpdateThreadBySlug(newData, slug)
}
func (s *ThreadService) UpdateThreadById(newData models.UpdateThread, id int) (models.Thread, models.Error) {
	// Check empty
	if newData.Message == "" && newData.Title == "" {
		return s.repo.GetThreadDataById(id)
	}

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

func (s *ThreadService) CreateThread(newThreadData models.Thread) (models.Thread, models.Error) {
	threadData, errCreate := s.repo.CreateThread(newThreadData)

	logrus.Println(errCreate)
	if errCreate.Code == http.StatusCreated {
		return threadData, errCreate
	}

	if errCreate.Code != http.StatusConflict {
		return threadData, errCreate
	}

	threadData, errGet := s.repo.GetThreadData(newThreadData.Slug)
	logrus.Println(errGet)

	if errGet.Code != http.StatusOK {
		return threadData, errGet
	}
	return threadData, errCreate
}
