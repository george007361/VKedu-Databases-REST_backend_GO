package postgres

import (
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
)

type ManagmentPostgres struct {
	db *sqlx.DB
}

var (
	queryClearDB = fmt.Sprintf("TRUNCATE %s, %s, %s, %s, %s CASCADE", userTable, forumTable, threadTable, voteTable, forumUsersTable)
)

func NewManagmentPostgres(db *sqlx.DB) *ManagmentPostgres {
	return &ManagmentPostgres{db: db}
}

func (r *ManagmentPostgres) Clear() models.Error {
	_, err := r.db.DB.Exec(queryClearDB)

	if err != nil {
		return models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return models.Error{Code: http.StatusOK, Message: "DB cleared succ"}
}

func (r *ManagmentPostgres) GetStatus() (models.Status, models.Error) {
	var status models.Status
	usersTableQuery := fmt.Sprintf("select count(*) from %s", userTable)
	forumTableQuery := fmt.Sprintf("select count(*) from %s", forumTable)
	threadTableQuery := fmt.Sprintf("select count(*) from %s", threadTable)
	postTableQuery := fmt.Sprintf("select count(*) from %s", postTable)

	err := r.db.DB.QueryRow(usersTableQuery).Scan(&status.User)
	if err != nil {
		return models.Status{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	err = r.db.DB.QueryRow(threadTableQuery).Scan(&status.Thread)
	if err != nil {
		return models.Status{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	err = r.db.DB.QueryRow(forumTableQuery).Scan(&status.Forum)
	if err != nil {
		return models.Status{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	err = r.db.DB.QueryRow(postTableQuery).Scan(&status.Post)
	if err != nil {
		return models.Status{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return status, models.Error{Code: http.StatusOK, Message: "Get status succ"}
}
