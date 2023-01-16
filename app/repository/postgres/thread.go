package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
)

type ThreadPostgres struct {
	db *sqlx.DB
}

var (
	querySelectThreadBySlug = fmt.Sprintf(`SELECT id, created, votes, title, author_nickname, forum_slug, message, slug
											FROM %s 
											WHERE slug=$1;`, threadTable)
	querySelectThreadById = fmt.Sprintf(`SELECT id, created, votes, title, author_nickname, forum_slug, message, slug
											FROM %s 
											WHERE id=$1;`, threadTable)
	queryCreateVote = fmt.Sprintf(`INSERT INTO %s (nickname, thread_id, voice) 
								  VALUES ($1, $2, $3) 
								  ON CONFLICT (nickname, thread_id) 
								  DO UPDATE SET voice=$3;`, voteTable)
	queryCreateThread = fmt.Sprintf(`INSERT INTO %s (title, author_nickname, message, forum_slug, slug, created) 
									VALUES ($1,$2,$3,$4,nullif($5,''),$6)
									RETURNING id, created, votes;`, threadTable)
)

func NewThreadPostgres(db *sqlx.DB) *ThreadPostgres {
	return &ThreadPostgres{db: db}
}

func (r *ThreadPostgres) CreateThread(newThreadData models.Thread) (models.Thread, models.Error) {
	err := r.db.DB.QueryRow(queryCreateThread, newThreadData.Title, newThreadData.AuthorNickname, newThreadData.Message, newThreadData.ForumSlug, newThreadData.Slug, newThreadData.Created).Scan(&newThreadData.ID, &newThreadData.Created, &newThreadData.Votes)

	if err != nil { // если такой форум уже еть
		fmt.Println(err)
		return models.Thread{}, models.Error{Code: 409, Message: err.Error()}
	}

	return newThreadData, models.Error{Code: http.StatusCreated, Message: "Thread created"}
}

func (r *ThreadPostgres) UpdateThreadBySlug(newData models.UpdateThread, slug string) (models.Thread, models.Error) {
	wherePart := "slug=$1"
	return r.updateThread(newData, wherePart, slug)
}

func (r *ThreadPostgres) UpdateThreadById(newData models.UpdateThread, id int) (models.Thread, models.Error) {
	wherePart := "id=$1"
	return r.updateThread(newData, wherePart, id)
}

func (r *ThreadPostgres) updateThread(newData models.UpdateThread, wherePart string, param interface{}) (models.Thread, models.Error) {
	queryParams := make([]interface{}, 0)
	queryParams = append(queryParams, param)

	setPart := ""
	if newData.Message != "" {
		queryParams = append(queryParams, newData.Message)
		setPart = setPart + fmt.Sprintf("message=$%d,", len(queryParams))
	}

	if newData.Title != "" {
		queryParams = append(queryParams, newData.Title)
		setPart = setPart + fmt.Sprintf("title=$%d,", len(queryParams))
	}

	if setPart != "" {
		setPart = setPart[:len(setPart)-1]
	}

	query := fmt.Sprintf(`UPDATE %s
						  SET %s
						  WHERE %s
						  RETURNING id, title, author_nickname, forum_slug, message, votes, slug, created;`, threadTable, setPart, wherePart)

	var threadData models.Thread
	err := r.db.DB.QueryRow(query, queryParams...).Scan(
		&threadData.ID,
		&threadData.Title,
		&threadData.AuthorNickname,
		&threadData.ForumSlug,
		&threadData.Message,
		&threadData.Votes,
		&threadData.Slug,
		&threadData.Created)
	if err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread "%v" not found`, param)}
	}

	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return threadData, models.Error{Code: http.StatusOK, Message: "Thread updated successfully"}
}

func (r *ThreadPostgres) GetThreadData(slug string) (models.Thread, models.Error) {

	var threadData models.Thread
	var threadSlug *string
	err := r.db.DB.QueryRow(querySelectThreadBySlug, slug).Scan(
		&threadData.ID,
		&threadData.Created,
		&threadData.Votes,
		&threadData.Title,
		&threadData.AuthorNickname,
		&threadData.ForumSlug,
		&threadData.Message,
		&threadSlug)
	if threadSlug != nil {
		threadData.Slug = *threadSlug
	}
	if err != nil && err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" was not found`, slug)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return threadData, models.Error{Code: http.StatusOK, Message: "Thread Data get succ"}
}

func (r *ThreadPostgres) GetThreadDataById(id int) (models.Thread, models.Error) {
	var threadData models.Thread

	var threadSlug *string

	err := r.db.DB.QueryRow(querySelectThreadById, id).Scan(
		&threadData.ID,
		&threadData.Created,
		&threadData.Votes,
		&threadData.Title,
		&threadData.AuthorNickname,
		&threadData.ForumSlug,
		&threadData.Message,
		&threadSlug)
	if threadSlug != nil {
		threadData.Slug = *threadSlug
	}

	if err != nil && err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" was not found`, id)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return threadData, models.Error{Code: http.StatusOK, Message: "Thread Data get succ"}
}

func (r *ThreadPostgres) VoteThread(vote models.Vote, id int) (models.Thread, models.Error) {

	_, err := r.db.DB.Exec(queryCreateVote, vote.Nickname, id, vote.Voice)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: err.Error()}
	}

	return r.GetThreadDataById(id)
}
