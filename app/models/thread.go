package models

import "time"

type Thread struct {
	Author  string    `json:"author" binding:"required"`  // Пользователь, создавший данную тему.
	Created time.Time `json:"created" `                   // Дата создания ветки на форуме.
	Forum   string    `json:"forum"`                      // Форум, в котором расположена данная ветка обсуждения.
	ID      int       `json:"id"`                         // Идентификатор ветки обсуждения.
	Message string    `json:"message" binding:"required"` // Описание ветки обсуждения.
	Title   string    `json:"title" binding:"required"`   // Заголовок ветки обсуждения.
	Slug    string    `json:"slug"`                       // Человекопонятный URL
	Votes   int       `json:"votes"`                      // Кол-во голосов непосредственно за данное сообщение форума.
}

type CreateThread struct {
	Author  string `json:"author" binding:"required"`  // Пользователь, создавший данную тему.
	Message string `json:"message" binding:"required"` // Описание ветки обсуждения.
	Title   string `json:"title" binding:"required"`   // Заголовок ветки обсуждения.
}

type UpdateThread struct {
	Message string `json:"message"` // Описание ветки обсуждения.
	Title   string `json:"title"`   // Заголовок ветки обсуждения.
}

type Vote struct {
	Nickname string `json:"nickname"` // Идентификатор пользователя.
	Voice    int32  `json:"voice"`    // Отданный голос.
}

type ThreadGetPostsParams struct {
	Limit int
	Since int
	Desc  bool
	Sort  string
}
