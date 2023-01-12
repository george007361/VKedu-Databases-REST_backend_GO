package repository

import (
	"github.com/george007361/db-course-proj/app/models"
	"github.com/george007361/db-course-proj/app/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type User interface {
	CreateUser(newUserData models.UserCreate) models.Error
	GetUserProfilesByEmailOrNickname(email string, nickname string) ([]*models.User, models.Error)
	GetUserProfile(nickname string) (models.User, models.Error)
	UpdateUserProfile(updatedData models.UserUpdate) (models.User, models.Error)
}

type Forum interface {
	CreateForum(newForumData models.Forum) (models.Forum, models.Error)
	GetForumData(slug string) (models.Forum, models.Error)
	GetForumUsers(params models.ForumUsersQueryParams) ([]models.User, models.Error)
	GetForumThreads(params models.ForumThreadsQueryParams) ([]models.Thread, models.Error)
}

type Thread interface {
	CreateThread(newThreadData models.Thread) (models.Thread, models.Error)
	GetThreadData(slug string) (models.Thread, models.Error)
	GetThreadDataById(id int) (models.Thread, models.Error)
	// CreatePostsByThreadSlug(newPostsData []models.Post, slug string) ([]models.Post, models.Error)
	// CreatePostsByThreadId(newPostsData []models.Post, id int) ([]models.Post, models.Error)
	UpdateThreadBySlug(newData models.UpdateThread, slug string) (models.Thread, models.Error)
	UpdateThreadById(newData models.UpdateThread, id int) (models.Thread, models.Error)
	GetThreadPostsBySlug(params models.ThreadGetPostsParams, slug string) ([]models.Post, models.Error)
	GetThreadPostsById(params models.ThreadGetPostsParams, id int) ([]models.Post, models.Error)
	VoteThreadBySlug(vote models.Vote, slug string) (models.Thread, models.Error)
	VoteThreadById(vote models.Vote, id int) (models.Thread, models.Error)
	//
	CreatePosts(newPostsData []models.Post, threadId int, forumSlug string) ([]models.Post, models.Error)
}

type Post interface {
	GetPostData(id int) (models.Post, models.Error)
	UpdatePostData(newData models.PostUpdate, id int) (models.Post, models.Error)
}

type Managment interface {
	Clear() models.Error
	GetStatus() (models.Status, models.Error)
}

type Repository struct {
	User
	Forum
	Thread
	Post
	Managment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:      postgres.NewUserPostgres(db),
		Forum:     postgres.NewForumPostgres(db),
		Thread:    postgres.NewThreadPostgres(db),
		Post:      postgres.NewPostPostgres(db),
		Managment: postgres.NewManagmentPostgres(db),
	}
}
