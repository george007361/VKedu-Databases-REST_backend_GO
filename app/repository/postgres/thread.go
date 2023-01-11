package postgres

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

func (r *ThreadPostgres) CreatePostsByThreadId(newPostsData []models.Post, id int) ([]models.Post, models.Error) {
	// Given : [parent, author, message]
	// Thread id

	// Check thread exists
	// Check forum exists
	checkThreadQuery := fmt.Sprintf(`select id, forum from %s where id=$1;`, threadTable)
	var threadId int
	var forumSlug string
	err := r.db.DB.QueryRow(checkThreadQuery, id).Scan(&threadId, &forumSlug)
	if err != nil && err == sql.ErrNoRows {
		// Если не нашёл thread по slug
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, id)}
	}

	logrus.Println(threadId)
	logrus.Println(forumSlug)

	return r.createPosts(newPostsData, threadId, forumSlug)
}

func (r *ThreadPostgres) CreatePostsByThreadSlug(newPostsData []models.Post, slug string) ([]models.Post, models.Error) {
	// Given : [parent, author, message]
	// Thread Slug

	// Check thread exists
	// Check forum exists
	checkThreadQuery := fmt.Sprintf(`select id, forum from %s where slug=$1;`, threadTable)
	var threadId int
	var forumSlug string
	err := r.db.DB.QueryRow(checkThreadQuery, slug).Scan(&threadId, &forumSlug)
	if err != nil && err == sql.ErrNoRows {
		// Если не нашёл thread по slug
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, slug)}
	}

	logrus.Println(threadId)
	logrus.Println(forumSlug)
	return r.createPosts(newPostsData, threadId, forumSlug)
}

func (r *ThreadPostgres) createPosts(newPostsData []models.Post, threadId int, forumSlug string) ([]models.Post, models.Error) {

	creationTime := time.Now()

	// Validate data
	for _, post := range newPostsData {
		// Check User exists
		checkUserQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
		err := r.db.DB.QueryRow(checkUserQuery, post.Author).Scan(&post.Author)
		if err != nil && err == sql.ErrNoRows {
			return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, post.Author)}
		}

		// Check parent exists
		if post.Parent != 0 {
			checkParentQuery := fmt.Sprintf(`select id from %s where id=$1;`, postTable)
			err = r.db.DB.QueryRow(checkParentQuery, post.Parent).Scan(&post.Parent)
			if err != nil && err == sql.ErrNoRows {
				return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Cant create post with parent id="%d". Parent not found`, post.Parent)}
			}
		}
	}

	// Upload posts

	// Fill values
	valuesString := ""
	valuesQueryParams := make([]interface{}, 0)
	cntQParams := 1
	for idx, _ := range newPostsData {
		newPostsData[idx].Thread = threadId
		newPostsData[idx].Forum = forumSlug
		newPostsData[idx].Created = creationTime
		valuesString += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d),`, cntQParams, cntQParams+1, cntQParams+2, cntQParams+3, cntQParams+4, cntQParams+5)
		valuesQueryParams = append(valuesQueryParams, newPostsData[idx].Thread, newPostsData[idx].Author, newPostsData[idx].Forum, newPostsData[idx].Message, newPostsData[idx].Parent, newPostsData[idx].Created)
		cntQParams += 6
	}
	// Trim last ','
	logrus.Println("------------------------------  ", valuesString)
	valuesString = valuesString[:len(valuesString)-1]

	logrus.Println(valuesString)
	logrus.Println(valuesQueryParams...)

	// Upload
	uploadQuery := fmt.Sprintf(`insert into %s (thread, author, forum, message, parent, created) values %s returning id;`, postTable, valuesString)
	rows, err := r.db.DB.Query(uploadQuery, valuesQueryParams...)

	if err != nil {
		return []models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	idx := 0
	for rows.Next() {
		err := rows.Scan(&newPostsData[idx].ID)
		if err != nil {
			return []models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		idx += 1
	}

	return newPostsData, models.Error{Code: http.StatusCreated, Message: "Posts created"}
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
	query := fmt.Sprintf(`update %s
						  set title=$2, message=$3
						  where %s
						  returning id, title, author, forum, message, votes, slug, created;`, threadTable, wherePart)
	var threadData models.Thread
	err := r.db.DB.QueryRow(query, param, newData.Title, newData.Message).Scan(
		&threadData.ID,
		&threadData.Title,
		&threadData.Author,
		&threadData.Forum,
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

func (r *ThreadPostgres) GetThreadPostsBySlug(params models.ThreadGetPostsParams, slug string) ([]models.Post, models.Error) {
	checkQuery := fmt.Sprintf(`select id from %s where slug=$1`, threadTable)
	var id int
	err := r.db.DB.QueryRow(checkQuery, slug).Scan(&id)
	if err == sql.ErrNoRows {
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, slug)}
	}
	return r.getThreadPosts(params, id)
}

func (r *ThreadPostgres) GetThreadPostsById(params models.ThreadGetPostsParams, id int) ([]models.Post, models.Error) {
	checkQuery := fmt.Sprintf(`select id from %s where id=$1`, threadTable)
	err := r.db.DB.QueryRow(checkQuery, id).Scan(&id)
	if err == sql.ErrNoRows {
		return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, id)}
	}
	return r.getThreadPosts(params, id)
}

func (r *ThreadPostgres) getThreadPosts(params models.ThreadGetPostsParams, id int) ([]models.Post, models.Error) {
	logrus.Println(params)
	var posts []models.Post
	whereStatement := ""
	orderStatement := ""
	limitStatement := ""
	queryParams := make([]interface{}, 0)

	switch params.Sort {
	case "flat":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = " thread = $1 AND id < $2 "
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " created desc "
			} else {
				whereStatement = " thread = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " created desc "
			}
		} else {
			if params.Since != 0 {
				whereStatement = " thread = $1 AND id > $2 "
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " created "
			} else {
				whereStatement = " thread = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " created "
			}
		}
	case "tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread = $1 and path < (select path from %s where id=$2) ", postTable)
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path desc "
			} else {
				whereStatement = " thread = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path desc "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread = $1 and path (select path from %s where id=$2) ", postTable)
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path "
			} else {
				whereStatement = " thread = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path "
			}
		}
	case "parent_tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path[1] in (select id from %s where parent = 0 and thread=$1 and id < (secect path[1] from %s where id = $2) order id desc limit $3) ", postTable, postTable)
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path[1] desc, path "
			} else {

				whereStatement = fmt.Sprintf(" path[1] in (select id from %s where parent=0 and thread=$1 order by id desc limit $2)", postTable)
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path[1] desc, path "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path[1] in (select id from %s where parent=0 and thread=$1 and id > (select path[1] from %s where id=$2) order by id limit $3) ", postTable, postTable)
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path "
			} else {

				whereStatement = fmt.Sprintf(" path[1] in (select id from %s where parent=0 and thread=$1 order by id limit $2) ", postTable)
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path "
			}
		}
	default:
		return posts, models.Error{Code: http.StatusBadRequest, Message: "Unknown sorting type"}
	}

	query := fmt.Sprintf(`select id, parent, author, message, isedited, forum, thread, created from %s
	where %s
	order by %s
	%s;`, postTable, whereStatement, orderStatement, limitStatement)

	rows, err := r.db.DB.Query(query, queryParams...)
	if err != nil {
		return []models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	for rows.Next() {
		post := models.Post{}
		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			fmt.Println("error while scan", err)
			return []models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		posts = append(posts, post)
	}

	return posts, models.Error{Code: http.StatusOK, Message: "posts succ"}
}

func (r *ThreadPostgres) VoteThreadBySlug(vote models.Vote, slug string) (models.Thread, models.Error) {
	checkQuery := fmt.Sprintf(`select id from %s where slug=$1`, threadTable)
	var id int
	err := r.db.DB.QueryRow(checkQuery, slug).Scan(&id)
	if err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with slug "%s" not found`, slug)}
	}
	return r.voteThread(vote, id)
}

func (r *ThreadPostgres) VoteThreadById(vote models.Vote, id int) (models.Thread, models.Error) {
	checkQuery := fmt.Sprintf(`select id from %s where id=$1`, threadTable)
	err := r.db.DB.QueryRow(checkQuery, id).Scan(&id)
	if err == sql.ErrNoRows {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Thread with id "%d" not found`, id)}
	}
	return r.voteThread(vote, id)
}

func (r *ThreadPostgres) voteThread(vote models.Vote, id int) (models.Thread, models.Error) {
	// Check user
	userQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
	err := r.db.DB.QueryRow(userQuery, vote.Nickname).Scan(&vote.Nickname)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, vote.Nickname)}
	}

	query := fmt.Sprintf(`insert into %s (nickname, thread, voice) values ($1, $2, $3) on conflict (nickname, thread) do update set voice=$3`, voteTable)

	_, err = r.db.DB.Exec(query, vote.Nickname, id, vote.Voice)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: err.Error()}
	}

	return r.GetThreadDataById(id)
}
