package models

import "time"

type Post struct {
	AuthorNickname string    `json:"author" binding:"required"`
	Created        time.Time `json:"created"`
	ForumSlug      string    `json:"forum"`
	ID             int       `json:"id"`
	Message        string    `json:"message" binding:"required"`
	ParentID       int       `json:"parent"`
	IsEdited       bool      `json:"isEdited"`
	ThreadId       int       `json:"thread"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type PostAllData struct {
	Post   Post    `json:"post"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
