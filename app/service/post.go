package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository"
)

type PostService struct {
	postRepo   repository.Post
	threadRepo repository.Thread
	userRepo   repository.User
}

func NewPostService(postRepo repository.Post, threadRepo repository.Thread, userRepo repository.User) *PostService {
	return &PostService{
		postRepo:   postRepo,
		threadRepo: threadRepo,
		userRepo:   userRepo,
	}
}

func (s *PostService) GetPostData(id int) (models.Post, models.Error) {
	return s.postRepo.GetPostData(id)
}

func (s *PostService) UpdatePostData(newData models.PostUpdate, id int) (models.Post, models.Error) {
	postData, err := s.postRepo.GetPostData(id)
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

	return s.postRepo.UpdatePostData(newData, id)
}

func (s *PostService) GetPostsByThreadSlug(params models.ThreadGetPostsParams, threadSlug string) ([]*models.Post, models.Error) {
	// Check thread exists
	threadData, err := s.threadRepo.GetThreadData(threadSlug)
	if err.Code != http.StatusOK {
		return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, threadSlug)}

	}
	return s.postRepo.GetPosts(params, threadData.ID)
}

func (s *PostService) GetPostsByThreadId(params models.ThreadGetPostsParams, threadID int) ([]*models.Post, models.Error) {
	// Check thread exists
	_, err := s.threadRepo.GetThreadDataById(threadID)
	if err.Code != http.StatusOK {
		return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, threadID)}
	}
	return s.postRepo.GetPosts(params, threadID)
}

func (s *PostService) CreatePostsByThreadSlug(newPostsData []*models.Post, threadSlug string) ([]*models.Post, models.Error) {
	// Given : [parent, author, message]
	// Thread Slug

	// Check thread exists
	threadData, err := s.threadRepo.GetThreadData(threadSlug)
	if err.Code == http.StatusNotFound {
		return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, threadSlug)}
	}

	return s.createPosts(newPostsData, threadData)
}

func (s *PostService) CreatePostsByThreadId(newPostsData []*models.Post, threadID int) ([]*models.Post, models.Error) {
	// Given : [parent, author, message]
	// Thread Id

	// Check thread exists
	threadData, err := s.threadRepo.GetThreadDataById(threadID)
	if err.Code == http.StatusNotFound {
		return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, threadID)}
	}

	return s.createPosts(newPostsData, threadData)
}

func (s *PostService) createPosts(newPostsData []*models.Post, threadData models.Thread) ([]*models.Post, models.Error) {
	creationTime := time.Now()
	// Check empty posts list
	if len(newPostsData) < 1 {
		return newPostsData, models.Error{Code: http.StatusCreated}
	}

	// Validate and fill
	for _, post := range newPostsData {
		// Check User exists
		userData, err := s.userRepo.GetUserProfile(post.AuthorNickname)
		if err.Code != http.StatusOK {
			return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, post.AuthorNickname)}
		}

		// Check parent exists and post thread eq parent thread
		if post.ParentID != 0 {
			parentData, err := s.postRepo.GetPostData(post.ParentID)
			// 1
			if err.Code == http.StatusNotFound {
				return nil, models.Error{Code: http.StatusConflict, Message: fmt.Sprintf(`Cant create post with parent id="%d". Parent not found`, post.ParentID)}
			}
			// 2
			if parentData.ThreadId != threadData.ID {
				return nil, models.Error{Code: http.StatusConflict, Message: "Parent post was created in another thread"}
			}

			if err.Code != http.StatusOK {
				return nil, err
			}
		}

		post.AuthorNickname = userData.Nickname
		post.ThreadId = threadData.ID
		post.ForumSlug = threadData.ForumSlug
		post.Created = creationTime
	}

	return s.postRepo.CreatePosts(newPostsData, threadData.ID, threadData.ForumSlug)
}
