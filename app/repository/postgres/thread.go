package postgres

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/george007361/db-course-proj/app/models"
	"github.com/jmoiron/sqlx"
)

type ThreadPostgres struct {
	db *sqlx.DB
}

func NewThreadPostgres(db *sqlx.DB) *ThreadPostgres {
	return &ThreadPostgres{db: db}
}

func (r *ThreadPostgres) GetThreadData(slug string) (models.Thread, models.Error) {
	query := fmt.Sprintf(` select id, created, votes, title, author_nickname, forum_slug, message, slug
							from %s
							where slug=$1`, threadTable)

	var threadData models.Thread
	var threadSlug *string
	err := r.db.DB.QueryRow(query, slug).Scan(
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
	query := fmt.Sprintf(` select id, created, votes, title, author_nickname, forum_slug, message, slug
							from %s
							where id=$1`, threadTable)

	var threadData models.Thread

	var threadSlug *string

	err := r.db.DB.QueryRow(query, id).Scan(
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

func (r *ThreadPostgres) CreatePosts(newPostsData []models.Post, threadId int, forumSlug string) ([]models.Post, models.Error) {

	creationTime := time.Now()

	// Validate data
	for _, post := range newPostsData {
		// Check User exists
		checkUserQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
		err := r.db.DB.QueryRow(checkUserQuery, post.AuthorNickname).Scan(&post.AuthorNickname)
		if err != nil && err == sql.ErrNoRows {
			return []models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, post.AuthorNickname)}
		}

		// Check parent exists and post thread eq parent thread
		if post.ParentID != 0 {
			checkParentQuery := fmt.Sprintf(`select id, thread_id from %s where id=$1;`, postTable)
			var parentThread int
			err = r.db.DB.QueryRow(checkParentQuery, post.ParentID).Scan(&post.ParentID, &parentThread)

			// 1
			if err != nil && err == sql.ErrNoRows {
				return []models.Post{}, models.Error{Code: http.StatusConflict, Message: fmt.Sprintf(`Cant create post with parent id="%d". Parent not found`, post.ParentID)}
			}
			// 2
			if parentThread != threadId {
				return []models.Post{}, models.Error{Code: http.StatusConflict, Message: "Parent post was created in another thread"}
			}
		}
	}

	// Upload posts

	// Fill values
	valuesString := ""
	valuesQueryParams := make([]interface{}, 0)
	cntQParams := 1
	for idx, _ := range newPostsData {
		newPostsData[idx].ThreadId = threadId
		newPostsData[idx].ForumSlug = forumSlug
		newPostsData[idx].Created = creationTime
		valuesString += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d),`, cntQParams, cntQParams+1, cntQParams+2, cntQParams+3, cntQParams+4, cntQParams+5)
		valuesQueryParams = append(valuesQueryParams, newPostsData[idx].ThreadId, newPostsData[idx].AuthorNickname, newPostsData[idx].ForumSlug, newPostsData[idx].Message, newPostsData[idx].ParentID, newPostsData[idx].Created)
		cntQParams += 6
	}
	// Trim last ','
	valuesString = valuesString[:len(valuesString)-1]

	// Upload
	uploadQuery := fmt.Sprintf(`insert into %s (thread_id, author_nickname, forum_slug, message, parent_id, created) values %s returning id, created;`, postTable, valuesString)
	rows, err := r.db.DB.Query(uploadQuery, valuesQueryParams...)

	if err != nil {
		return []models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	idx := 0
	for rows.Next() {
		err := rows.Scan(&newPostsData[idx].ID, &newPostsData[idx].Created)
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

	query := fmt.Sprintf(`update %s
						  set %s
						  where %s
						  returning id, title, author_nickname, forum_slug, message, votes, slug, created;`, threadTable, setPart, wherePart)

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
	var posts []models.Post
	whereStatement := ""
	orderStatement := ""
	limitStatement := ""
	queryParams := make([]interface{}, 0)

	switch params.Sort {
	case "flat":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = " thread_id = $1 AND id < $2 "
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " created desc, id desc "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " created desc, id desc"
			}
		} else {
			if params.Since != 0 {
				whereStatement = " thread_id = $1 AND id > $2 "
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " created, id "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " created, id "
			}
		}
	case "tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread_id = $1 and path_tree < (select path_tree from %s where id=$2) ", postTable)
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path_tree desc "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path_tree desc "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread_id = $1 and path_tree > (select path_tree from %s where id=$2) ", postTable)
				limitStatement = "limit $3"
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path_tree "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "limit $2"
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path_tree "
			}
		}
	case "parent_tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path_tree[1] in (select id from %s where parent_id = 0 and thread_id=$1 and id < (select path_tree[1] from %s where id = $2) order by id desc limit $3) ", postTable, postTable)
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path_tree[1] desc, path_tree "
			} else {

				whereStatement = fmt.Sprintf(" path_tree[1] in (select id from %s where parent_id=0 and thread_id=$1 order by id desc limit $2)", postTable)
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path_tree[1] desc, path_tree "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path_tree[1] in (select id from %s where parent_id=0 and thread_id=$1 and id > (select path_tree[1] from %s where id=$2) order by id limit $3) ", postTable, postTable)
				queryParams = append(queryParams, id, params.Since, params.Limit)
				orderStatement = " path_tree "
			} else {

				whereStatement = fmt.Sprintf(" path_tree[1] in (select id from %s where parent_id=0 and thread_id=$1 order by id limit $2) ", postTable)
				queryParams = append(queryParams, id, params.Limit)
				orderStatement = " path_tree "
			}
		}
	default:
		return posts, models.Error{Code: http.StatusBadRequest, Message: "Unknown sorting type"}
	}

	query := fmt.Sprintf(`select id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created from %s
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
			&post.ParentID,
			&post.AuthorNickname,
			&post.Message,
			&post.IsEdited,
			&post.ForumSlug,
			&post.ThreadId,
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

	query := fmt.Sprintf(`insert into %s (nickname, thread_id, voice) values ($1, $2, $3) on conflict (nickname, thread_id) do update set voice=$3`, voteTable)

	_, err = r.db.DB.Exec(query, vote.Nickname, id, vote.Voice)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: err.Error()}
	}

	return r.GetThreadDataById(id)
}

func (r *ThreadPostgres) CreateThread(newThreadData models.Thread) (models.Thread, models.Error) {

	userQuery := fmt.Sprintf(`select nickname from %s where nickname=$1;`, userTable)
	err := r.db.DB.QueryRow(userQuery, newThreadData.AuthorNickname).Scan(&newThreadData.AuthorNickname)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`User with nickname "%s" not found`, newThreadData.AuthorNickname)}
	}

	forumQuery := fmt.Sprintf(`select slug from %s where slug=$1;`, forumTable)
	err = r.db.DB.QueryRow(forumQuery, newThreadData.ForumSlug).Scan(&newThreadData.ForumSlug)
	if err != nil {
		return models.Thread{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Forum with slug "%s" not found`, newThreadData.AuthorNickname)}
	}

	threadQuery := fmt.Sprintf(`
	insert into %s 
    (title, author_nickname, message, forum_slug, slug, created) 
	values ($1,$2,$3,$4,nullif($5,''),$6) 
	returning id, created, votes`, threadTable)

	err = r.db.DB.QueryRow(threadQuery, newThreadData.Title, newThreadData.AuthorNickname, newThreadData.Message, newThreadData.ForumSlug, newThreadData.Slug, newThreadData.Created).Scan(&newThreadData.ID, &newThreadData.Created, &newThreadData.Votes)

	if err != nil { // если такой форум уже еть
		fmt.Println(err)
		return models.Thread{}, models.Error{Code: 409, Message: err.Error()}
	}

	return newThreadData, models.Error{Code: http.StatusCreated, Message: "Thread created"}
}
