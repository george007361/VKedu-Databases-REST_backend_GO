package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(newUserData models.UserCreate) models.Error {
	query := fmt.Sprintf(`insert into %s (nickname, fullname, about, email) values($1, $2, $3, $4);`, userTable)

	_, err := r.db.DB.Exec(query, newUserData.Nickname, newUserData.FullName, newUserData.About, newUserData.Email)
	if err != nil {
		return models.Error{Code: http.StatusConflict, Message: err.Error()}
	}

	return models.Error{Code: http.StatusCreated, Message: "User created"}
}

func (r *UserPostgres) GetUserProfilesByEmailOrNickname(email string, nickname string) ([]*models.User, models.Error) {
	query := fmt.Sprintf(`select nickname, fullname, about, email from %s where email = $1 or nickname = $2;`, userTable)

	rows, err := r.db.DB.Query(query, email, nickname)
	if err != nil {
		return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	usersData := make([]*models.User, 0)
	cnt := 0
	for rows.Next() {
		var userData models.User
		err := rows.Scan(&userData.Nickname, &userData.FullName, &userData.About, &userData.Email)
		if err != nil {
			return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		usersData = append(usersData, &userData)
		cnt += 1
	}

	if cnt == 0 {
		return nil, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Users with email "%s" or nickname "%s" not found`, email, nickname)}
	}

	return usersData, models.Error{Code: http.StatusOK, Message: "Get users success"}
}

func (r *UserPostgres) GetUserProfile(nickname string) (models.User, models.Error) {
	query := fmt.Sprintf(`select nickname, fullname, about, email from %s where nickname = $1 limit 1;`, userTable)
	var userData models.User

	err := r.db.DB.QueryRow(query, nickname).Scan(&userData.Nickname, &userData.FullName, &userData.About, &userData.Email)

	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, nickname)}
	}

	if err != nil {
		return models.User{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return userData, models.Error{Code: http.StatusOK, Message: "Get user success"}
}

func (r *UserPostgres) UpdateUserProfile(updatedData models.UserUpdate) (models.User, models.Error) {
	var userUpdatedData models.User

	query := fmt.Sprintf(`
	update %s set 
	fullname=coalesce(nullif($1, ''), fullname),
	about=coalesce(nullif($2, ''), about),
	email=coalesce(nullif($3, ''), email)
	where nickname=$4
	returning fullname, about, email, nickname`, userTable)

	err := r.db.DB.QueryRow(query, updatedData.FullName, updatedData.About, updatedData.Email, updatedData.Nickname).Scan(&userUpdatedData.FullName, &userUpdatedData.About, &userUpdatedData.Email, &userUpdatedData.Nickname)

	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, updatedData.Nickname)}
	}

	if err != nil {
		return models.User{}, models.Error{Code: 409, Message: err.Error()}
	}

	return userUpdatedData, models.Error{Code: http.StatusOK, Message: "User data updated successfully"}
}
