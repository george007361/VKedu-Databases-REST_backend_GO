package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ForumPostgres struct {
	db *sqlx.DB
}

func NewForumPostgres(db *sqlx.DB) *ForumPostgres {
	return &ForumPostgres{db: db}
}

func (r *ForumPostgres) CreateForum(newForumData models.Forum) (models.Forum, models.Error) {

	userQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
	err := r.db.DB.QueryRow(userQuery, newForumData.User).Scan(&newForumData.User)
	if err != nil {
		return models.Forum{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, newForumData.User)}
	}

	forumQuery := fmt.Sprintf(`insert into %s (slug, title, "user")	values ($1, $2, $3) returning slug, title, "user", posts, threads`, forumTable)
	var forumData models.Forum

	err = r.db.DB.QueryRow(forumQuery, newForumData.Slug, newForumData.Title, newForumData.User).Scan(&forumData.Slug, &forumData.Title, &forumData.User, &forumData.Posts, &forumData.Threads)
	logrus.Println(err)

	if err != nil && err != sql.ErrNoRows { // если такой форум уже еcть
		return models.Forum{}, models.Error{Code: http.StatusConflict, Message: "Forum already exists"}
	}

	return forumData, models.Error{Code: http.StatusCreated, Message: "Forum created"}
}

func (r *ForumPostgres) GetForumData(slug string) (models.Forum, models.Error) {
	query := fmt.Sprintf(`select slug, title, "user", posts, threads from %s where slug=$1`, forumTable)

	var forumData models.Forum

	err := r.db.DB.QueryRow(query, slug).Scan(&forumData.Slug, &forumData.Title, &forumData.User, &forumData.Posts, &forumData.Threads)

	if err != nil && err == sql.ErrNoRows {
		return models.Forum{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, slug)}
	}

	if err != nil {
		return models.Forum{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return forumData, models.Error{Code: http.StatusOK, Message: "Forum data get succ"}
}

func (r *ForumPostgres) GetForumUsers(params models.ForumUsersQueryParams) ([]models.User, models.Error) {

	checkQuery := fmt.Sprintf(`select slug from %s where slug=$1;`, forumTable)
	err := r.db.DB.QueryRow(checkQuery, params.Slug).Scan(&params.Slug)
	if err != nil && err == sql.ErrNoRows {
		return []models.User{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, params.Slug)}
	}

	var queryParams []interface{}
	queryParams = append(queryParams, params.Slug, params.Limit)
	whereStatementStr := ""
	orderStatementStr := ""

	if params.Since != "" {
		if params.Desc {
			whereStatementStr = " and nickname < $3"
		} else {
			whereStatementStr = " and nickname > $3"
		}
		queryParams = append(queryParams, params.Since)
	}

	if params.Desc {
		orderStatementStr = "desc"
	}

	query := fmt.Sprintf(`select nickname, fullname, about, email from %s
							where forum=$1 %s
							order by nickname %s 
							limit $2;`,
		forumUsersTable, whereStatementStr, orderStatementStr)

	rows, err := r.db.DB.Query(query, queryParams...)

	if err != nil {
		return []models.User{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	users := make([]models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email,
		)
		logrus.Println(err)

		if err != nil {
			return []models.User{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		users = append(users, *user)
	}
	return users, models.Error{Code: http.StatusOK, Message: "Forum users list get succ"}
}

func (r *ForumPostgres) GetForumThreads(params models.ForumThreadsQueryParams) ([]models.Thread, models.Error) {

	checkQuery := fmt.Sprintf(`select slug from %s where slug=$1;`, forumTable)
	err := r.db.DB.QueryRow(checkQuery, params.Slug).Scan(&params.Slug)
	if err != nil && err == sql.ErrNoRows {
		return []models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, params.Slug)}
	}

	var queryParams []interface{}
	queryParams = append(queryParams, params.Slug, params.Limit)
	whereStatementStr := ""
	orderStatementStr := ""

	if params.Since != -1 {
		if params.Desc {
			whereStatementStr = " and extract(EPOCH from created)  <= $3 "
		} else {
			whereStatementStr = " and extract(EPOCH from created)  >= $3 "
		}
		queryParams = append(queryParams, params.Since)
	}

	if params.Desc {
		orderStatementStr = "desc"
	}

	query := fmt.Sprintf(`select id, slug, forum, author, title, message, votes, created 
						from %s 
						where forum = $1 %s
						order by created %s
						limit $2`,
		threadTable, whereStatementStr, orderStatementStr)

	rows, err := r.db.DB.Query(query, queryParams...)
	if err != nil {
		return []models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	threads := make([]models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		err = rows.Scan(
			&thread.ID,
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		if err != nil {
			return []models.Thread{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		threads = append(threads, *thread)
	}
	return threads, models.Error{Code: http.StatusOK, Message: "Forum threads list get succ"}
}

func (r *ForumPostgres) CreateThreadInForum(newThreadData models.Thread) (models.Thread, models.Error) {

	userQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
	err := r.db.DB.QueryRow(userQuery, newThreadData.Author).Scan(&newThreadData.Author)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, newThreadData.Author)}
	}

	forumQuery := fmt.Sprintf(`select slug from %s where slug=$1;`, forumTable)
	err = r.db.DB.QueryRow(forumQuery, newThreadData.Forum).Scan(&newThreadData.Forum)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, newThreadData.Forum)}
	}

	threadQuery := fmt.Sprintf(`
	insert into %s 
    (title, author, message, forum, slug) 
	values ($1,$2,$3,$4,$5) 
	returning id, created, votes`, threadTable)

	err = r.db.DB.QueryRow(threadQuery, newThreadData.Title, newThreadData.Author, newThreadData.Message, newThreadData.Forum, newThreadData.Slug).Scan(&newThreadData.ID, &newThreadData.Created, &newThreadData.Votes)

	if err != nil { // если такой форум уже еть
		fmt.Println(err)
		return models.Thread{}, models.Error{Code: 409, Message: err.Error()}
	}

	return newThreadData, models.Error{Code: http.StatusCreated, Message: "Thread created"}
}
