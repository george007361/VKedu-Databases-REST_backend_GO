package models

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Status struct {
	User   int32 `json:"user"`   // Кол-во пользователей в базе данных
	Forum  int32 `json:"forum"`  // Кол-во разделов в базе данных
	Thread int32 `json:"thread"` // Кол-во веток обсуждения в базе данных
	Post   int64 `json:"post"`   // Кол-во сообщений в базе данных
}
