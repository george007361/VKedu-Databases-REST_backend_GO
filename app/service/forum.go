package service

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type ForumService struct {
	repo repository.Forum
}

func NewForumService(repo repository.Forum) *ForumService {
	return &ForumService{repo: repo}
}

func (s *ForumService) CreateForum(newForumData models.Forum) (models.Forum, models.Error) {
	return s.repo.CreateForum(newForumData)
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

func (s *ForumService) CreateThreadInForum(newThreadData models.Thread) (models.Thread, models.Error) {
	return s.repo.CreateThreadInForum(newThreadData)
}

// func (s *UserService) GetUserProfile(nickname string) (models.User, models.Error) {
// 	return s.repo.GetUserProfile(nickname)
// }

// func (s *UserService) UpdateUserProfile(userData models.User) (models.User, models.Error) {
// 	return s.repo.UpdateUserProfile(userData)
// }
