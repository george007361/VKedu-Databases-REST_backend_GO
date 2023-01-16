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

var (
	querySelectPostByID = fmt.Sprintf(`SELECT id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created
										FROM %s
										WHERE id=$1;`, postTable)
	queryUpdatePostByID = fmt.Sprintf(`UPDATE %s
										SET message=$1, isedited=true 
										WHERE id=$2 
										RETURNING id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created;`, postTable)
	querySelectPostIdThreadIdByID = fmt.Sprintf(`SELECT id, thread_id FROM %s WHERE id=$1;`, postTable)
)

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{db: db}
}

func (r *PostPostgres) GetPostData(id int) (models.Post, models.Error) {

	var postData models.Post

	err := r.db.DB.QueryRow(querySelectPostByID, id).Scan(
		&postData.ID,
		&postData.ParentID,
		&postData.AuthorNickname,
		&postData.Message,
		&postData.IsEdited,
		&postData.ForumSlug,
		&postData.ThreadId,
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

	var postData models.Post
	err := r.db.DB.QueryRow(queryUpdatePostByID, newData.Message, id).Scan(
		&postData.ID,
		&postData.ParentID,
		&postData.AuthorNickname,
		&postData.Message,
		&postData.IsEdited,
		&postData.ForumSlug,
		&postData.ThreadId,
		&postData.Created)

	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf(`Post with id="%d" not found`, id)}
	}

	if err != nil {
		return models.Post{}, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return postData, models.Error{Code: http.StatusOK, Message: "Post successfully updated"}
}

func (r *PostPostgres) GetPosts(params models.ThreadGetPostsParams, threadID int) ([]*models.Post, models.Error) {
	whereStatement := ""
	orderStatement := ""
	limitStatement := ""
	queryParams := make([]interface{}, 0)

	switch params.Sort {
	case "flat":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = " thread_id = $1 AND id < $2 "
				limitStatement = "LIMIT $3"
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " created DESC, id DESC "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "LIMIT $2"
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " created DESC, id DESC"
			}
		} else {
			if params.Since != 0 {
				whereStatement = " thread_id = $1 AND id > $2 "
				limitStatement = "LIMIT $3"
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " created, id "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "LIMIT $2"
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " created, id "
			}
		}
	case "tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread_id = $1 AND path_tree < (SELECT path_tree FROM %s WHERE id=$2) ", postTable)
				limitStatement = "LIMIT $3"
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " path_tree DESC "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "LIMIT $2"
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " path_tree DESC "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" thread_id = $1 AND path_tree > (SELECT path_tree FROM %s WHERE id=$2) ", postTable)
				limitStatement = "LIMIT $3"
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " path_tree "
			} else {
				whereStatement = " thread_id = $1 "
				limitStatement = "LIMIT $2"
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " path_tree "
			}
		}
	case "parent_tree":
		if params.Desc {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path_tree[1] IN (SELECT id FROM %s WHERE parent_id = 0 AND thread_id=$1 AND id < (SELECT path_tree[1] FROM %s WHERE id = $2) ORDER BY id DESC LIMIT $3) ", postTable, postTable)
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " path_tree[1] DESC, path_tree "
			} else {

				whereStatement = fmt.Sprintf(" path_tree[1] IN (SELECT id FROM %s WHERE parent_id=0 AND thread_id=$1 ORDER BY id DESC LIMIT $2)", postTable)
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " path_tree[1] DESC, path_tree "
			}
		} else {
			if params.Since != 0 {
				whereStatement = fmt.Sprintf(" path_tree[1] IN (SELECT id FROM %s WHERE parent_id=0 AND thread_id=$1 AND id > (SELECT path_tree[1] FROM %s WHERE id=$2) ORDER BY id LIMIT $3) ", postTable, postTable)
				queryParams = append(queryParams, threadID, params.Since, params.Limit)
				orderStatement = " path_tree "
			} else {

				whereStatement = fmt.Sprintf(" path_tree[1] IN (SELECT id FROM %s WHERE parent_id=0 AND thread_id=$1 ORDER BY id LIMIT $2) ", postTable)
				queryParams = append(queryParams, threadID, params.Limit)
				orderStatement = " path_tree "
			}
		}
	default:
		return nil, models.Error{Code: http.StatusBadRequest, Message: "Unknown sorting type"}
	}

	query := fmt.Sprintf(`SELECT id, parent_id, author_nickname, message, isedited, forum_slug, thread_id, created
						FROM %s 
						WHERE %s 
						ORDER BY %s	
						%s;`, postTable, whereStatement, orderStatement, limitStatement)

	rows, err := r.db.DB.Query(query, queryParams...)

	if err != nil {
		return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	posts := make([]*models.Post, 0)
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
			return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		posts = append(posts, &post)
	}
	return posts, models.Error{Code: http.StatusOK, Message: "posts succ"}
}

func (r *PostPostgres) CreatePosts(newPostsData []*models.Post, threadId int, forumSlug string) ([]*models.Post, models.Error) {
	// Fill values
	valuesString := ""
	valuesQueryParams := make([]interface{}, 0)
	cntQParams := 1
	for idx, _ := range newPostsData {
		valuesString += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d, $%d),`, cntQParams, cntQParams+1, cntQParams+2, cntQParams+3, cntQParams+4, cntQParams+5)
		valuesQueryParams = append(valuesQueryParams, newPostsData[idx].ThreadId, newPostsData[idx].AuthorNickname, newPostsData[idx].ForumSlug, newPostsData[idx].Message, newPostsData[idx].ParentID, newPostsData[idx].Created)
		cntQParams += 6
	}
	// Trim last ','
	valuesString = valuesString[:len(valuesString)-1]

	// Upload
	uploadQuery := fmt.Sprintf(`INSERT INTO %s (thread_id, author_nickname, forum_slug, message, parent_id, created)
								VALUES %s
								RETURNING id, created;`, postTable, valuesString)
	rows, err := r.db.DB.Query(uploadQuery, valuesQueryParams...)

	if err != nil {
		return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	idx := 0
	for rows.Next() {
		err := rows.Scan(&newPostsData[idx].ID, &newPostsData[idx].Created)
		if err != nil {
			return nil, models.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		idx += 1
	}

	return newPostsData, models.Error{Code: http.StatusCreated, Message: "Posts created"}
}
