package models

type User struct {
	Email    string `json:"email" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	About    string `json:"about" binding:"required"`
}

type UserCreate struct {
	Email    string `json:"email" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	About    string `json:"about" binding:"required"`
}

type UserUpdate struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Nickname string `json:"nickname"`
	About    string `json:"about"`
}
