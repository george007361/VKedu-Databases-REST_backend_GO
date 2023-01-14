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
	query := fmt.Sprintf(` select id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created
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
		return models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Post with id "%d" was not found`, id)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return postData, models.Error{Code: http.StatusOK, Message: "Post Data get succ"}
}

func (r *PostPostgres) UpdatePostData(newData models.PostUpdate, id int) (models.Post, models.Error) {
	query := fmt.Sprintf(`update %s
						  set message=$1, isedited=true 
						  where id=$2 
						  returning id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created`, postTable)

	var postData models.Post
	err := r.db.DB.QueryRow(query, newData.Message, id).Scan(
		&postData.ID,
		&postData.Parent,
		&postData.Author,
		&postData.Message,
		&postData.IsEdited,
		&postData.Forum,
		&postData.Thread,
		&postData.Created)

	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Post with id="%d" not found`, id)}
	}

	if err != nil {
		return models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return postData, models.Error{Code: http.StatusOK, Message: "Post successfully updated"}
}
