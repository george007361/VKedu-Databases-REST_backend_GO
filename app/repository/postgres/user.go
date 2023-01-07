package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user models.User) models.Error {
	query := fmt.Sprintf(`insert into %s (nickname, fullname, about, email) values($1, $2, $3, $4);`, userTable)

	_, err := r.db.DB.Exec(query, user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		return models.Error{Code: http.StatusConflict, Message: err.Error()}
	}

	return models.Error{Code: http.StatusCreated, Message: "User created"}
}

func (r *UserPostgres) GetUserProfile(nickname string) (models.User, models.Error) {
	query := fmt.Sprintf(`select nickname, fullname, about, email from %s where nickname = $1 limit 1;`, userTable)
	var userData models.User

	err := r.db.DB.QueryRow(query, nickname).Scan(&userData.Nickname, &userData.FullName, &userData.About, &userData.Email)

	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, nickname)}
	}

	logrus.Printf("%v+\n", userData)
	return userData, models.Error{Code: http.StatusOK, Message: "Get user success"}
}

func (r *UserPostgres) UpdateUserProfile(userData models.User) (models.User, models.Error) {
	var userUpdatedData models.User

	query := fmt.Sprintf(`
	update %s set 
	fullname=coalesce(nullif($1, ''), fullname),
	about=coalesce(nullif($2, ''), about),
	email=coalesce(nullif($3, ''), email)
	where nickname=$4
	returning fullname, about, email, nickname`, userTable)

	err := r.db.DB.QueryRow(query, userData.FullName, userData.About, userData.Email, userData.Nickname).Scan(&userUpdatedData.FullName, &userUpdatedData.About, &userUpdatedData.Email, &userUpdatedData.Nickname)

	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, userData.Nickname)}
	}

	if err != nil {
		return models.User{}, models.Error{Code: 409, Message: err.Error()}
	}

	return userUpdatedData, models.Error{Code: http.StatusOK, Message: "User data updated successfully"}
}
