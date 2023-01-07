package models

import "time"

type Post struct { //  Сообщение внутри ветки обсуждения на форуме.
	Author   string    `json:"author" binding:"required"`  // Автор, написавший данное сообщение.
	Created  time.Time `json:"created"`                    // Дата создания сообщения на форуме.
	Forum    string    `json:"forum"`                      // Идентификатор форума (slug) данного сообещния.
	ID       int       `json:"id"`                         // Идентификатор данного сообщения.
	Message  string    `json:"message" binding:"required"` // Собственно сообщение форума.
	Parent   int       `json:"parent" binding:"required"`  // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	IsEdited bool      `json:"isEdited"`                   // Истина, если данное сообщение было изменено.
	Thread   int       `json:"thread"`                     // Идентификатор форума (slug) данного сообещния.
}

type PostUpdate struct {
	Message string `json:"message"` // Собственно сообщение форума.
}

type PostAllData struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}
