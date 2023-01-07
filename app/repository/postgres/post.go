package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
)

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{db: db}
}

func (r *PostPostgres) GetPostData(id int) (models.Post, models.Error) {
	query := fmt.Sprintf(` select id, parent, author, message, isedited, forum, thread, created
							from %s
							where id=$1`, postTable)

	var postData models.Post

	err := r.db.DB.QueryRow(query, id).Scan(
		&postData.ID,
		&postData.Parent,
		&postData.Author,
		&postData.Message,
		&postData.IsEdited,
		&postData.Forum,
		&postData.Thread,
		&postData.Created)

	if err != nil && err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprint(`Post with id "%d" was not found`, id)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return postData, models.Error{Code: http.StatusOK, Message: "Post Data get succ"}
}
