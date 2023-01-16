package service

import (
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type ThreadService struct {
	threadRepo repository.Thread
	userRepo   repository.User
	forumRepo  repository.Forum
}

func NewThreadService(threadRepo repository.Thread, userRepo repository.User, forumRepo repository.Forum) *ThreadService {
	return &ThreadService{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (s *ThreadService) GetThreadData(slug string) (models.Thread, models.Error) {
	return s.threadRepo.GetThreadData(slug)
}

func (s *ThreadService) GetThreadDataById(id int) (models.Thread, models.Error) {
	return s.threadRepo.GetThreadDataById(id)
}

func (s *ThreadService) UpdateThreadBySlug(newData models.UpdateThread, slug string) (models.Thread, models.Error) {
	// Check empty
	if newData.Message == "" && newData.Title == "" {
		return s.threadRepo.GetThreadData(slug)
	}

	return s.threadRepo.UpdateThreadBySlug(newData, slug)
}
func (s *ThreadService) UpdateThreadById(newData models.UpdateThread, id int) (models.Thread, models.Error) {
	// Check empty
	if newData.Message == "" && newData.Title == "" {
		return s.threadRepo.GetThreadDataById(id)
	}

	return s.threadRepo.UpdateThreadById(newData, id)
}

func (s *ThreadService) VoteThreadBySlug(vote models.Vote, slug string) (models.Thread, models.Error) {
	// Check user
	userData, err := s.userRepo.GetUserProfile(vote.Nickname)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, vote.Nickname)}
	}
	vote.Nickname = userData.Nickname

	// Check thread
	threadData, err := s.threadRepo.GetThreadData(slug)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, slug)}
	}

	return s.threadRepo.VoteThread(vote, threadData.ID)
}

func (s *ThreadService) VoteThreadById(vote models.Vote, id int) (models.Thread, models.Error) {
	// Check user
	userData, err := s.userRepo.GetUserProfile(vote.Nickname)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, vote.Nickname)}
	}
	vote.Nickname = userData.Nickname

	// Check thread
	_, err = s.threadRepo.GetThreadDataById(id)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, id)}
	}
	return s.threadRepo.VoteThread(vote, id)
}

func (s *ThreadService) CreateThread(newThreadData models.Thread) (models.Thread, models.Error) {

	// Check user
	userData, err := s.userRepo.GetUserProfile(newThreadData.AuthorNickname)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, newThreadData.AuthorNickname)}
	}
	newThreadData.AuthorNickname = userData.Nickname

	// Check forum
	forumData, err := s.forumRepo.GetForumData(newThreadData.ForumSlug)
	if err.Code != http.StatusOK {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, newThreadData.ForumSlug)}
	}
	newThreadData.ForumSlug = forumData.Slug

	threadData, errCreate := s.threadRepo.CreateThread(newThreadData)
	if errCreate.Code == http.StatusCreated {
		return threadData, errCreate
	}

	if errCreate.Code != http.StatusConflict {
		return threadData, errCreate
	}

	threadData, errGet := s.threadRepo.GetThreadData(newThreadData.Slug)

	if errGet.Code != http.StatusOK {
		return threadData, errGet
	}
	return threadData, errCreate
}
