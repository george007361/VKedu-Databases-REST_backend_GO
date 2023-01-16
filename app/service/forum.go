package service

import (
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type ForumService struct {
	forumRepo repository.Forum
	userRepo  repository.User
}

func NewForumService(forumRepo repository.Forum, userRepo repository.User) *ForumService {
	return &ForumService{
		forumRepo: forumRepo,
		userRepo:  userRepo,
	}
}

func (s *ForumService) CreateForum(newForumData models.Forum) (models.Forum, models.Error) {

	// Check user
	userData, err := s.userRepo.GetUserProfile(newForumData.AuthorNickname)
	if err.Code != http.StatusOK {
		return models.Forum{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, newForumData.AuthorNickname)}
	}
	newForumData.AuthorNickname = userData.Nickname

	// Try to create
	forumData, errCreate := s.forumRepo.CreateForum(newForumData)
	if errCreate.Code != http.StatusConflict {
		return forumData, errCreate
	}

	// Get existing forum data
	forumData, errGet := s.forumRepo.GetForumData(newForumData.Slug)
	if errGet.Code != http.StatusOK {
		return models.Forum{}, errGet
	}

	return forumData, errCreate
}

func (s *ForumService) GetForumData(slug string) (models.Forum, models.Error) {
	return s.forumRepo.GetForumData(slug)
}

func (s *ForumService) GetForumUsers(params models.ForumUsersQueryParams) ([]models.User, models.Error) {

	// Check forum exists
	_, err := s.forumRepo.GetForumData(params.Slug)
	if err.Code != http.StatusOK {
		return []models.User{}, err
	}

	return s.forumRepo.GetForumUsers(params)
}

func (s *ForumService) GetForumThreads(params models.ForumThreadsQueryParams) ([]models.Thread, models.Error) {
	// Check forum exists
	_, err := s.forumRepo.GetForumData(params.Slug)
	if err.Code != http.StatusOK {
		return []models.Thread{}, err
	}

	return s.forumRepo.GetForumThreads(params)
}
