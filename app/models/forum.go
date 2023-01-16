package models

type Forum struct {
	Slug           string `json:"slug" binding:"required"`
	Title          string `json:"title" binding:"required"`
	AuthorNickname string `json:"user" binding:"required"`
	Posts          int    `json:"posts"`
	Threads        int    `json:"threads"`
}

type ForumUsersQueryParams struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
}

type ForumThreadsQueryParams struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
}
