package service

import (
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
	"github.com/sirupsen/logrus"
)

type ForumService struct {
	repo repository.Forum
}

func NewForumService(repo repository.Forum) *ForumService {
	return &ForumService{repo: repo}
}

func (s *ForumService) CreateForum(newForumData models.Forum) (models.Forum, models.Error) {
	forumData, errCreate := s.repo.CreateForum(newForumData)
	if errCreate.Code != http.StatusConflict {
		return forumData, errCreate
	}

	forumData, errGet := s.repo.GetForumData(newForumData.Slug)
	if errGet.Code != http.StatusOK {
		return forumData, errGet
	}
	logrus.Println(forumData, errCreate, errGet)
	return forumData, errCreate
}

func (s *ForumService) GetForumData(slug string) (models.Forum, models.Error) {
	return s.repo.GetForumData(slug)
}

func (s *ForumService) GetForumUsers(params models.ForumUsersQueryParams) ([]models.User, models.Error) {
	return s.repo.GetForumUsers(params)
}

func (s *ForumService) GetForumThreads(params models.ForumThreadsQueryParams) ([]models.Thread, models.Error) {
	return s.repo.GetForumThreads(params)
}
