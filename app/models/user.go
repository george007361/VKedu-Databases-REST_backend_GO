package models

type User struct {
	Email    string `json:"email" binding:"required"`    // Почтовый адрес пользователя (уникальное поле)
	FullName string `json:"fullname" binding:"required"` // Полное имя пользователя
	Nickname string `json:"nickname" binding:"required"` //Имя пользователя (уникальное поле)
	About    string `json:"about" binding:"required"`    // Описание пользователя
}

type UserCreate struct {
	Email    string `json:"email" binding:"required"`    // Почтовый адрес пользователя (уникальное поле)
	FullName string `json:"fullname" binding:"required"` // Полное имя пользователя
	Nickname string `json:"nickname" binding:"required"` //Имя пользователя (уникальное поле)
	About    string `json:"about" binding:"required"`    // Описание пользователя
}

type UserUpdate struct {
	Email    string `json:"email"`    // Почтовый адрес пользователя (уникальное поле)
	FullName string `json:"fullname"` // Полное имя пользователя
	Nickname string `json:"nickname"` //Имя пользователя (уникальное поле)
	About    string `json:"about"`    // Описание пользователя
}
