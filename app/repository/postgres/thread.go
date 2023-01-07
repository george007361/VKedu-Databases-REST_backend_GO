package postgres

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type ThreadPostgres struct {
	db *sqlx.DB
}

func NewThreadPostgres(db *sqlx.DB) *ThreadPostgres {
	return &ThreadPostgres{db: db}
}

func (r *ThreadPostgres) GetThreadData(slug string) (models.Thread, models.Error) {
	query := fmt.Sprintf(` select id, created, votes, title, author, forum, message, slug
							from %s
							where slug=$1`, threadTable)

	var threadData models.Thread

	err := r.db.DB.QueryRow(query, slug).Scan(
		&threadData.ID,
		&threadData.Created,
		&threadData.Votes,
		&threadData.Title,
		&threadData.Author,
		&threadData.Forum,
		&threadData.Message,
		&threadData.Slug)
	if err != nil && err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprint(`Thread with slug "%s" was not found`, slug)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return threadData, models.Error{Code: http.StatusOK, Message: "Thread Data get succ"}
}

func (r *ThreadPostgres) GetThreadDataById(id int) (models.Thread, models.Error) {
	query := fmt.Sprintf(` select id, created, votes, title, author, forum, message, slug
							from %s
							where id=$1`, threadTable)

	var threadData models.Thread

	err := r.db.DB.QueryRow(query, id).Scan(
		&threadData.ID,
		&threadData.Created,
		&threadData.Votes,
		&threadData.Title,
		&threadData.Author,
		&threadData.Forum,
		&threadData.Message,
		&threadData.Slug)
	if err != nil && err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprint(`Thread with id "%d" was not found`, id)}
	}

	if err != nil && err != sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return threadData, models.Error{Code: http.StatusOK, Message: "Thread Data get succ"}
}

func (r *ThreadPostgres) CreatePostsByThreadSlug(newPostsData []models.Post, slug string) ([]models.Post, models.Error) {
	if len(newPostsData) != 0 && newPostsData[0].Parent != 0 {
		var parentCheck int
		checkQuery := fmt.Sprintf("select thread from %s where id = $1", postTable)
		err := r.db.DB.QueryRow(checkQuery, newPostsData[0].Parent).Scan(&parentCheck)
		if err != nil {
			return nil, models.Error{Code: http.StatusConflict, Message: "Parent post was created in another thread"}
		}
	}

	creationTime := time.Now()

	query := fmt.Sprintf(`insert into %s (thread, author, forum, message, parent, created) values `, postTable)

	queryParams := make([]interface{}, 0)

	last := len(newPostsData) - 1
	for i, post := range newPostsData {

		if !db.CheckValidParent(thread.ID, post.Parent) {
			return nil, models.Error{Code: 409}
		}

		var parentPost *int
		if post.Parent != 0 {
			parentPost = new(int)
			*parentPost = post.Parent
		}

		if i == last {
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d) `, i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		} else {
			//query += fmt.Sprintf(`(%d,'%s','%s','%s', %d, $1), `, parentPost, post.Author, post.Message, thread.Forum, thread.ID)
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d), `, i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		}
		queryParams = append(queryParams, thread.ID, post.Author, thread.Forum, post.Message, parentPost, createdTime)
	}

	query += " returning id, created"

	//transaction, err := db.DB.Begin()
	//rows, err := transaction.Query(query, queryParams...)
	//if err != nil {
	//	transaction.Rollback()
	//	_, ok := err.(pgx.PgError)
	//	if ok {
	//		return nil, models.Error{Code: 404, Message: fmt.Sprintf("%d", 1)}
	//	}
	//}

	//transaction, err := db.DbCreate.Begin(context.Background())
	//batch := new(pgx.Batch)
	//batch.Queue(query, queryParams...)
	rows, err := db.DB.Query(query, queryParams...)

	fmt.Println("ADDING POST ERROR ", err)
	if err != nil {
		return nil, models.Error{Code: 500}
	}

	i := 0
	for rows.Next() {
		err = rows.Scan(&posts[i].ID, &posts[i].Created)
		if err != nil {
			return nil, models.Error{Code: 500}
		}
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		fmt.Println(posts[i])
		i++
	}

	if dbErr, ok := rows.Err().(pgx.PgError); ok {
		fmt.Println("PGX ERROR")
		switch dbErr.Code {
		case pgerrcode.RaiseException:
			fmt.Println("40404")
			return nil, models.Error{Code: 404, Message: "Post parent not found"}
		case "23503":
			fmt.Println("23503")
			return nil, models.Error{Code: 404, Message: "User not found"}
		}
	}

	return posts, models.Error{}
}
