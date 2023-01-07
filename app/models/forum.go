package models

type Forum struct {
	Slug    string `json:"slug" binding:"required"`  // Человекопонятный URL
	Title   string `json:"title" binding:"required"` // Название форума
	User    string `json:"user" binding:"required"`  // Nickname пользователя, который отвечает за форум
	Posts   int    `json:"posts"`                    // Общее кол-во сообщений в данном форуме
	Threads int    `json:"threads"`                  // Общее кол-во ветвей обсуждения в данном форуме
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
	Since int
	Desc  bool
}
