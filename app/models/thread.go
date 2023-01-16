package models

import "time"

type Thread struct {
	AuthorNickname string    `json:"author" binding:"required"`
	Created        time.Time `json:"created" `
	ForumSlug      string    `json:"forum"`
	ID             int       `json:"id"`
	Message        string    `json:"message" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Slug           string    `json:"slug"`
	Votes          int       `json:"votes"`
}

type CreateThread struct {
	AuthorNickname string `json:"author" binding:"required"`
	Message        string `json:"message" binding:"required"`
	Title          string `json:"title" binding:"required"`
}

type UpdateThread struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}

type ThreadGetPostsParams struct {
	Limit int
	Since int
	Desc  bool
	Sort  string
}
